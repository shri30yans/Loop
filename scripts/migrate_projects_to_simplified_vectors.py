#!/usr/bin/env python3
"""
Script to migrate existing projects to the simplified vector storage system.
This script:
1. Fetches all projects from PostgreSQL
2. For each project, extracts entities and generates a summary
3. Generates embeddings for the summaries
4. Stores the embeddings in the simplified project_vectors table (only project_id, summary, embedding)
5. Stores metadata (technologies, domains, stakeholders, features) in Neo4j
"""

import os
import sys
import json
import psycopg2
import requests
from psycopg2.extras import RealDictCursor
import argparse
import time
from neo4j import GraphDatabase

# Configuration
DEFAULT_EMBEDDING_SERVICE_URL = "http://localhost:5000"
DEFAULT_BATCH_SIZE = 10

def connect_to_db(db_url):
    """Connect to PostgreSQL database"""
    try:
        conn = psycopg2.connect(db_url)
        return conn
    except Exception as e:
        print(f"Error connecting to database: {e}")
        sys.exit(1)

def connect_to_neo4j(uri, user, password):
    """Connect to Neo4j database"""
    try:
        driver = GraphDatabase.driver(uri, auth=(user, password))
        return driver
    except Exception as e:
        print(f"Error connecting to Neo4j: {e}")
        sys.exit(1)

def fetch_projects(conn, batch_size=DEFAULT_BATCH_SIZE, offset=0):
    """Fetch projects from database in batches"""
    cursor = conn.cursor(cursor_factory=RealDictCursor)
    cursor.execute(
        """
        SELECT p.project_id, p.title, p.description, 
               json_agg(ps.content) as sections
        FROM projects p
        LEFT JOIN project_sections ps ON p.project_id = ps.project_id
        GROUP BY p.project_id, p.title, p.description
        ORDER BY p.created_at
        LIMIT %s OFFSET %s
        """,
        (batch_size, offset)
    )
    projects = cursor.fetchall()
    cursor.close()
    return projects

def count_projects(conn):
    """Count total number of projects"""
    cursor = conn.cursor()
    cursor.execute("SELECT COUNT(*) FROM projects")
    count = cursor.fetchone()[0]
    cursor.close()
    return count

def generate_summary(project, ai_api_url):
    """Generate summary for a project using AI API"""
    # Combine project sections
    sections = project.get('sections', [])
    if sections and sections[0] is not None:
        combined_text = "\n\n".join(sections)
    else:
        combined_text = project.get('description', '')
    
    # Prepare prompt
    prompt = f"""
    Please analyze this project and provide:
    1. A concise summary (3-5 sentences)
    2. Key technologies used
    3. Main domain or industry
    4. Target stakeholders
    5. Key features

    Format your response as JSON with the following structure:
    {{
        "summary": "...",
        "technologies": ["tech1", "tech2", ...],
        "domains": ["domain1", "domain2", ...],
        "stakeholders": ["stakeholder1", "stakeholder2", ...],
        "features": ["feature1", "feature2", ...]
    }}

    Project Title: {project['title']}
    Project Description: {project['description']}
    Project Content:
    {combined_text}
    """
    
    # Call AI API
    try:
        response = requests.post(
            ai_api_url,
            json={"prompt": prompt, "max_tokens": 1000}
        )
        response.raise_for_status()
        result = response.json()
        
        # Parse the response to extract JSON
        content = result.get('choices', [{}])[0].get('text', '').strip()
        
        # Find JSON in the response
        start_idx = content.find('{')
        end_idx = content.rfind('}') + 1
        if start_idx >= 0 and end_idx > start_idx:
            json_str = content[start_idx:end_idx]
            try:
                return json.loads(json_str)
            except json.JSONDecodeError:
                print(f"Failed to parse JSON from response: {json_str}")
        
        print(f"Could not extract valid JSON from response: {content}")
        return {
            "summary": project['description'],
            "technologies": [],
            "domains": [],
            "stakeholders": [],
            "features": []
        }
    except Exception as e:
        print(f"Error calling AI API: {e}")
        return {
            "summary": project['description'],
            "technologies": [],
            "domains": [],
            "stakeholders": [],
            "features": []
        }

def generate_embedding(text, embedding_service_url):
    """Generate embedding for text using embedding service"""
    try:
        response = requests.post(
            f"{embedding_service_url}/embed",
            json={"text": text}
        )
        response.raise_for_status()
        result = response.json()
        return result.get('embedding', [])
    except Exception as e:
        print(f"Error generating embedding: {e}")
        return []

def store_vector(conn, project_id, summary, embedding):
    """Store vector in simplified database schema"""
    cursor = conn.cursor()
    try:
        # Format embedding as a string for pgvector
        if embedding:
            vector_str = f"[{','.join(str(x) for x in embedding)}]"
        else:
            vector_str = "[]"
            
        cursor.execute(
            """
            INSERT INTO project_vectors
            (project_id, summary, embedding)
            VALUES (%s, %s, %s::vector)
            ON CONFLICT (project_id) 
            DO UPDATE SET
                summary = EXCLUDED.summary,
                embedding = EXCLUDED.embedding,
                updated_at = NOW()
            """,
            (
                project_id,
                summary,
                vector_str
            )
        )
        conn.commit()
        return True
    except Exception as e:
        conn.rollback()
        print(f"Error storing vector: {e}")
        return False
    finally:
        cursor.close()

def store_metadata_in_neo4j(neo4j_driver, project_id, project_title, metadata):
    """Store project metadata in Neo4j"""
    with neo4j_driver.session() as session:
        try:
            # Create project node if it doesn't exist
            session.run(
                """
                MERGE (p:Project {id: $project_id})
                SET p.name = $project_title
                """,
                project_id=str(project_id),
                project_title=project_title
            )
            
            # Add technology nodes and relationships
            for tech in metadata.get('technologies', []):
                session.run(
                    """
                    MERGE (t:Technology {name: $tech_name})
                    MERGE (p:Project {id: $project_id})
                    MERGE (p)-[:USES_TECHNOLOGY]->(t)
                    """,
                    tech_name=tech,
                    project_id=str(project_id)
                )
            
            # Add domain nodes and relationships
            for domain in metadata.get('domains', []):
                session.run(
                    """
                    MERGE (d:Domain {name: $domain_name})
                    MERGE (p:Project {id: $project_id})
                    MERGE (p)-[:IN_DOMAIN]->(d)
                    """,
                    domain_name=domain,
                    project_id=str(project_id)
                )
            
            # Add stakeholder nodes and relationships
            for stakeholder in metadata.get('stakeholders', []):
                session.run(
                    """
                    MERGE (s:Stakeholder {name: $stakeholder_name})
                    MERGE (p:Project {id: $project_id})
                    MERGE (p)-[:HAS_STAKEHOLDER]->(s)
                    """,
                    stakeholder_name=stakeholder,
                    project_id=str(project_id)
                )
            
            # Add feature nodes and relationships
            for feature in metadata.get('features', []):
                session.run(
                    """
                    MERGE (f:Feature {name: $feature_name})
                    MERGE (p:Project {id: $project_id})
                    MERGE (p)-[:HAS_FEATURE]->(f)
                    """,
                    feature_name=feature,
                    project_id=str(project_id)
                )
            
            return True
        except Exception as e:
            print(f"Error storing metadata in Neo4j: {e}")
            return False

def main():
    parser = argparse.ArgumentParser(description='Migrate projects to simplified vector storage')
    parser.add_argument('--db-url', required=True, help='PostgreSQL connection URL')
    parser.add_argument('--ai-api-url', required=True, help='AI API URL for summary generation')
    parser.add_argument('--embedding-service-url', default=DEFAULT_EMBEDDING_SERVICE_URL, 
                        help='Embedding service URL')
    parser.add_argument('--neo4j-uri', required=True, help='Neo4j URI')
    parser.add_argument('--neo4j-user', required=True, help='Neo4j username')
    parser.add_argument('--neo4j-password', required=True, help='Neo4j password')
    parser.add_argument('--batch-size', type=int, default=DEFAULT_BATCH_SIZE, 
                        help='Number of projects to process in each batch')
    args = parser.parse_args()
    
    # Connect to databases
    pg_conn = connect_to_db(args.db_url)
    neo4j_driver = connect_to_neo4j(args.neo4j_uri, args.neo4j_user, args.neo4j_password)
    
    # Count total projects
    total_projects = count_projects(pg_conn)
    print(f"Found {total_projects} projects to process")
    
    # Process projects in batches
    processed = 0
    offset = 0
    
    while processed < total_projects:
        # Fetch batch of projects
        projects = fetch_projects(pg_conn, args.batch_size, offset)
        if not projects:
            break
        
        print(f"Processing batch of {len(projects)} projects (offset: {offset})")
        
        for project in projects:
            project_id = project['project_id']
            print(f"Processing project {project_id}: {project['title']}")
            
            # Generate summary and metadata
            summary_data = generate_summary(project, args.ai_api_url)
            
            # Generate embedding for summary only
            embedding = generate_embedding(summary_data.get('summary', ''), args.embedding_service_url)
            
            # Store vector in PostgreSQL (simplified schema)
            if embedding:
                success = store_vector(pg_conn, project_id, summary_data.get('summary', ''), embedding)
                if success:
                    print(f"Successfully stored vector for project {project_id}")
                else:
                    print(f"Failed to store vector for project {project_id}")
            else:
                print(f"Failed to generate embedding for project {project_id}")
            
            # Store metadata in Neo4j
            success = store_metadata_in_neo4j(
                neo4j_driver, 
                project_id, 
                project['title'], 
                summary_data
            )
            if success:
                print(f"Successfully stored metadata for project {project_id}")
            else:
                print(f"Failed to store metadata for project {project_id}")
            
            # Sleep briefly to avoid overwhelming services
            time.sleep(0.5)
        
        processed += len(projects)
        offset += args.batch_size
        print(f"Processed {processed}/{total_projects} projects")
    
    pg_conn.close()
    neo4j_driver.close()
    print("Migration completed")

if __name__ == "__main__":
    main()
