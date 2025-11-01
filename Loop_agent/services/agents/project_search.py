"""
Project Search Agent implementation.
This agent is responsible for processing project search queries and coordinating with the Go backend.
"""

import logging
import re
from typing import Dict, List, Any, Optional, Tuple
from services.agents.base import BaseAgent

logger = logging.getLogger(__name__)

class ProjectSearchAgent(BaseAgent):
    """Agent responsible for processing project search queries"""
    
    def __init__(self, go_client):
        """
        Initialize the project search agent.
        
        Args:
            go_client: Client for communicating with the Go backend
        """
        super().__init__(
            name="Project Search Agent",
            role="Extract and process search queries for software projects"
        )
        self.go_client = go_client
        self.system_prompt = """
        You are an expert project discovery specialist with deep knowledge of software architecture.
        Your task is to analyze search queries, extract key terms, and determine the best search strategy.
        
        Search Strategy Guidelines:
        
        1. Vector Search (Similarity-Based):
           - Use for finding semantically similar projects
           - Best for 'like X' or 'similar to X' queries
           - Handles natural language and conceptual matching
           - Example: 'Show me projects similar to Netflix'
        
        2. Graph Search (Conditional):
           - Use for specific attribute-based filtering
           - Best for queries with explicit conditions
           - Handles relationships and technical requirements
           - Example: 'Find React projects using GraphQL'
        
        3. Hybrid Search (Two-Stage):
           - First: Vector search for relevant candidates
           - Then: Graph search to filter and augment results
           - Best for complex queries combining similarity and conditions
           - Example: 'Projects like Airbnb but using Vue.js'
        """
    
    async def extract_search_terms(self, query: str, context: List[Dict] = None) -> Dict:
        """
        Extract search terms and determine search mode.
        
        Args:
            query: The user's search query
            context: Previous conversation history
            
        Returns:
            Dict containing extracted terms and analysis
        """
        # Format conversation context if available
        context_str = ""
        if context and len(context) > 0:
            context_str = "Previous conversation:\n"
            for msg in context:
                role = msg.get("role", "unknown")
                content = msg.get("content", "")
                context_str += f"{role}: {content}\n"
        
        prompt = f"""Analyze this search query: "{query}"

{context_str}

Your task is to extract the core search terms and determine the best search mode.

For example:
- "Show me projects about 3D printing" → "3d printing" (vector search)
- "Find React projects using GraphQL" → "react graphql" (graph search)
- "Projects like Netflix but using Vue.js" → "netflix vue.js" (hybrid search)

Respond in JSON format:
{{
  "extracted_terms": "The core search terms, cleaned of filler words",
  "search_mode": "vector/graph/hybrid",
  "reasoning": "Brief explanation of your extraction and mode selection"
}}
"""
        
        response = await self.call_llm(self.system_prompt, prompt)
        analysis = self.parse_json_response(response)
        
        # Handle error cases
        if "error" in analysis:
            logger.error(f"Error extracting search terms: {analysis['error']}")
            # Return a safe default that will still allow search to proceed
            return {
                "extracted_terms": query,
                "search_mode": "vector",
                "reasoning": "Defaulting to vector search due to term extraction error"
            }
        
        # Validate the analysis has the required fields
        if not analysis.get("extracted_terms"):
            analysis["extracted_terms"] = query
        if not analysis.get("search_mode"):
            analysis["search_mode"] = "vector"
        
        logger.info(f"Extracted terms from '{query}': {analysis['extracted_terms']}")
        return analysis
    
    # Regex patterns for vague query detection
    VAGUE_QUERY_PATTERNS = {
        # Single word or very short queries
        r'^\s*\w{1,3}\s*$': 'Query too short',
        
        # Common greetings and conversation starters
        r'^\s*(hi|hello|hey|howdy|sup|yo|hiya|hola|greetings)\s*$': 'Greeting only',
        
        # Generic requests without specifics
        r'^\s*(show|find|get|give|list|display)\s+(me\s+)?(some\s+)?(projects?|repos?|apps?|code)\s*$': 'Generic request',
        
        # Questions without specific context
        r'^\s*(what|which|who|how).*(projects?|repos?|apps?)\s*\??\s*$': 'Question without context',
        
        # Single technology without context
        r'^\s*(using|with|in)\s+\w+\s*$': 'Single technology without context',
        
        # Very generic tech requests
        r'^\s*(good|best|top|popular|trending)\s+(projects?|repos?|apps?)\s*$': 'Generic quality request'
    }
    
    # Technical terms that indicate specificity
    TECHNICAL_TERMS = {
        # Programming languages
        'python', 'javascript', 'typescript', 'java', 'c++', 'c#', 'ruby', 'go', 'rust', 'php', 'swift',
        'kotlin', 'scala', 'perl', 'haskell', 'lua', 'r', 'dart', 'elixir', 'clojure', 'erlang',
        
        # Frameworks and libraries
        'react', 'angular', 'vue', 'svelte', 'django', 'flask', 'express', 'spring', 'laravel', 'rails',
        'tensorflow', 'pytorch', 'pandas', 'numpy', 'scikit', 'bootstrap', 'tailwind', 'jquery',
        'node', 'deno', 'next.js', 'nuxt', 'gatsby', 'fastapi', 'graphql', 'apollo', 'redux', 'mobx',
        
        # Databases
        'sql', 'mysql', 'postgresql', 'mongodb', 'firebase', 'dynamodb', 'redis', 'cassandra', 'neo4j',
        'sqlite', 'oracle', 'mariadb', 'couchdb', 'elasticsearch', 'supabase',
        
        # Cloud and DevOps
        'aws', 'azure', 'gcp', 'docker', 'kubernetes', 'terraform', 'jenkins', 'gitlab', 'github',
        'bitbucket', 'circleci', 'travis', 'ansible', 'puppet', 'chef', 'prometheus', 'grafana',
        
        # Mobile
        'android', 'ios', 'flutter', 'react-native', 'xamarin', 'ionic', 'cordova', 'swift', 'kotlin',
        
        # Domains and applications
        'ecommerce', 'blog', 'cms', 'authentication', 'payment', 'streaming', 'chat', 'messaging',
        'social', 'analytics', 'dashboard', 'booking', 'reservation', 'mapping', 'geolocation',
        'machine learning', 'ai', 'blockchain', 'cryptocurrency', 'iot', 'ar', 'vr', 'game',
        
        # Architecture patterns
        'microservices', 'serverless', 'spa', 'pwa', 'rest', 'graphql', 'mvc', 'mvvm', 'event-driven',
        'monolith', 'jamstack', 'headless', 'microfrontend'
    }
    
    # Topic words that are somewhat meaningful but generic
    TOPIC_WORDS = {
        # Project-related terms
        'project', 'projects', 'repository', 'repositories', 'repo', 'repos',
        'codebase', 'codebases', 'application', 'applications', 'app', 'apps',
        'system', 'systems', 'program', 'programs', 'software', 'code',

        # Technology generics
        'technology', 'technologies', 'tech', 'framework', 'frameworks',
        'library', 'libraries', 'tool', 'tools', 'platform', 'platforms',
        'language', 'languages', 'solution', 'solutions', 'service', 'services'
    }
    
    # Filler words that add little meaning
    FILLER_WORDS = {
        # Common query starters
        'show', 'find', 'get', 'give', 'list', 'display', 'search', 'query',
        'looking', 'need', 'want', 'seeking', 'tell', 'about', 'for', 'me',
        
        # Common prepositions and articles
        'the', 'a', 'an', 'of', 'in', 'on', 'with', 'by', 'to', 'from', 
        'at', 'for', 'about', 'as', 'into', 'like', 'through', 'after',
        'over', 'between', 'out', 'against', 'during', 'without', 'before',
        'under', 'around', 'among',
        
        # Common conjunctions
        'and', 'or', 'but', 'so', 'yet', 'nor', 'for', 'because', 'although',
        'since', 'unless', 'while', 'where', 'if', 'then', 'than',
        
        # Question words
        'what', 'which', 'who', 'whom', 'whose', 'when', 'where', 'why',
        'how', 'is', 'are', 'was', 'were', 'be', 'being', 'been', 'do',
        'does', 'did', 'has', 'have', 'had', 'can', 'could', 'will',
        'would', 'shall', 'should', 'may', 'might', 'must',
        
        # Common adjectives
        'good', 'best', 'better', 'nice', 'great', 'awesome', 'cool',
        'interesting', 'useful', 'helpful', 'popular', 'top', 'latest',
        'new', 'modern', 'simple', 'easy', 'complex', 'advanced'
    }
    
    def _extract_terms(self, query: str) -> Dict[str, List[str]]:
        """
        Extract different types of terms from the query.
        
        Args:
            query: The user's search query
            
        Returns:
            Dict containing categorized terms
        """
        words = query.lower().split()
        return {
            'technical': [word for word in words if word in self.TECHNICAL_TERMS],
            'topic': [word for word in words if word in self.TOPIC_WORDS],
            'filler': [word for word in words if word in self.FILLER_WORDS],
            'other': [word for word in words if word not in (
                self.TECHNICAL_TERMS | self.TOPIC_WORDS | self.FILLER_WORDS
            )]
        }
    
    def _calculate_query_score(self, query: str) -> float:
        """
        Calculate a score for the query based on its specificity.
        
        Args:
            query: The user's search query
            
        Returns:
            Float score between 0 and 1, higher is better
        """
        # Extract terms
        terms = self._extract_terms(query)
        
        # Count words by category
        word_count = len(query.split())
        if word_count == 0:
            return 0.0
            
        # Calculate weighted score
        weights = {
            'technical': 2.0,  # Technical terms are most valuable
            'other': 1.2,      # Other specific terms 
            'topic': 0.7,      # Topic words add some value
            'filler': 0.3      # Filler words add little value
        }
        
        weighted_sum = sum(
            len(terms[category]) * weights[category]
            for category in weights.keys()
        )
        
        # Normalize by word count with a slight bonus for longer queries
        length_factor = min(1.0, word_count / 5)  # Reaches 1.0 at 5+ words
        base_score = weighted_sum / max(1, word_count)
        
        return base_score * (0.7 + 0.3 * length_factor)  # Length provides up to 30% boost
    
    def _get_improvement_suggestions(self, query: str, terms: Dict[str, List[str]]) -> List[str]:
        """
        Generate specific suggestions to improve the query.
        
        Args:
            query: The user's search query
            terms: Extracted terms from the query
            
        Returns:
            List of suggestion strings
        """
        suggestions = []
        
        # Suggest adding technical terms if none present
        if not terms['technical']:
            suggestions.append("Add specific technologies or frameworks (e.g., Python, React, TensorFlow)")
        
        # Suggest adding domain if query is just about technology
        if terms['technical'] and not any(domain in query.lower() for domain in 
                                         ['ecommerce', 'blog', 'social', 'game', 'chat', 'analytics']):
            suggestions.append("Include a domain or application type (e.g., e-commerce, blog, social network)")
        
        # Suggest adding functionality if not present
        if not any(func in query.lower() for func in 
                  ['authentication', 'payment', 'search', 'messaging', 'dashboard']):
            suggestions.append("Describe specific functionality (e.g., authentication, data visualization)")
        
        # If query is very short, suggest expanding it
        if len(query.split()) < 3:
            suggestions.append("Expand your query with more details about what you're looking for")
        
        # If no specific suggestions were generated, add generic ones
        if not suggestions:
            suggestions = [
                "Mention specific technologies or frameworks",
                "Include domains or industries",
                "Describe functionality or features",
                "Add any specific requirements"
            ]
            
        return suggestions
    
    def is_vague_query(self, query: str) -> Tuple[bool, Dict]:
        """
        Check if a query is too vague to process meaningfully using regex patterns
        and weighted scoring.
        
        Args:
            query: The user's search query
            
        Returns:
            Tuple of (is_vague: bool, feedback: dict)
        """
        # Strip whitespace and convert to lowercase
        cleaned = query.strip().lower()
        
        # Check for very short queries
        if len(cleaned) < 3:
            return True, {
                'reason': 'Query too short',
                'suggestions': ["Please provide a more detailed query"]
            }
        
        # Check regex patterns first for common vague patterns
        for pattern, reason in self.VAGUE_QUERY_PATTERNS.items():
            if re.match(pattern, cleaned):
                logger.info(f"Query matched vague pattern '{pattern}': '{query}'")
                return True, {
                    'reason': reason,
                    'suggestions': ["Please provide a more specific query"]
                }
        
        # Extract terms and calculate score
        terms = self._extract_terms(cleaned)
        score = self._calculate_query_score(cleaned)
        
        # Generate feedback
        feedback = {
            'score': score,
            'terms': terms,
            'suggestions': self._get_improvement_suggestions(cleaned, terms)
        }
        
        # Determine if query is vague based on score
        is_vague = score < 0.6
        
        if is_vague:
            logger.info(f"Query is too vague (score: {score:.2f}): '{query}'")
            
        return is_vague, feedback
    
    async def search(self, query: str, context: List[Dict] = None, access_token: str = None) -> Dict:
        """
        Process a search query and return results.
        
        Args:
            query: The user's search query
            context: Previous conversation history
            access_token: JWT token for authentication
            
        Returns:
            Dict containing search results or improvement request
        """
        # Check if query is too vague
        is_vague, feedback = self.is_vague_query(query)
        if is_vague:
            # Create examples based on the original query to make suggestions more relevant
            examples = []
            
            # Try to extract any meaningful terms from the query
            terms = feedback.get('terms', {})
            tech_terms = terms.get('technical', [])
            other_terms = terms.get('other', [])
            
            # Generate examples based on extracted terms or defaults
            if tech_terms:
                # If we have technical terms, use them in examples
                tech = tech_terms[0]
                examples.append(f"Instead of just '{tech}', try '{tech} projects for user authentication'")
            elif other_terms:
                # If we have other specific terms, use them
                term = other_terms[0]
                examples.append(f"Instead of just '{term}', try '{term} applications with React and Node.js'")
            else:
                # Default examples
                examples = [
                    "Instead of just 'projects', try 'machine learning projects for image classification'",
                    "Instead of 'about tech', try 'React e-commerce projects with shopping cart'",
                    "Instead of 'hi', try 'Java backend projects with database integration'"
                ]
            
            # Format the message with specific suggestions from feedback
            suggestions = "\n".join(f"{i+1}. {suggestion}" for i, suggestion in enumerate(feedback.get('suggestions', [])))
            examples_text = "\n".join(f"- {example}" for example in examples[:2])
            
            return {
                "type": "results",
                "query": query,
                "extracted_terms": query,
                "search_mode": "none",
                "results": [],
                "message": f"""I notice your query is quite general. To help you find the most relevant projects, try:

{suggestions}

For example:
{examples_text}
"""
            }
        
        # Extract search terms with context
        analysis = await self.extract_search_terms(query, context)
        
        search_terms = analysis.get("extracted_terms", query)
        search_mode = analysis.get("search_mode", "vector")
        
        logger.info(f"Executing search - Terms: {search_terms}, Mode: {search_mode}")
        
        # Execute search via Go backend
        try:
            results = await self.go_client.search_tool(search_terms, mode=search_mode, access_token=access_token)
            return {
                "query": query,
                "extracted_terms": search_terms,
                "search_mode": search_mode,
                "results": results,
                "reasoning": analysis.get("reasoning", "")
            }
        except Exception as e:
            error_msg = f"Search error: {str(e)}"
            self.log_error(error_msg, {
                "query": query,
                "extracted_terms": search_terms,
                "search_mode": search_mode
            })
            return {
                "error": error_msg,
                "query": query,
                "extracted_terms": search_terms,
                "search_mode": search_mode
            }
    
    def format_results(self, search_results: Dict) -> str:
        """
        Format search results for display.
        
        Args:
            search_results: Dict containing search results and metadata
            
        Returns:
            Formatted string for display
        """
        query = search_results.get("query", "")
        extracted_terms = search_results.get("extracted_terms", "")
        search_mode = search_results.get("search_mode", "vector")
        results = search_results.get("results", [])
        
        if not results:
            return (
                "# No matching projects found\n\n"
                f"I couldn't find any projects matching your query: '{query}'\n\n"
                "## Suggestions to improve your search:\n\n"
                "1. Try more general terms.\n"
                "2. Check your spelling.\n"
                "3. Use alternative terminology.\n"
                "4. Explore related areas.\n"
                "5. Try a different search mode.\n\n"
                "Would you like to try a modified search query?"
            )
        
        # Start with the header and query analysis
        formatted_response = f"""# Search Results 🎯

## Query Analysis
I analyzed your query: "{query}"
Extracted search terms: "{extracted_terms}"

This query is best suited for {search_mode} search because {search_results.get('reasoning', '')}"""
        
        # Add primary matches section
        formatted_response += "\n\n## Primary Matches\n"
        
        for i, result in enumerate(results[:5], 1):
            project_name = result.get("name", f"Project {i}")
            description = result.get("description", "No description available")
            score = result.get("score", 0)
            
            formatted_response += f"\n### {i}. {project_name}\n"
            formatted_response += f"**Relevance Score**: {score:.2f}\n\n"
            formatted_response += f"{description}\n"
            
            # Add technical details if available
            if "technologies" in result:
                techs = result.get("technologies", [])
                if techs:
                    formatted_response += "\n**Technologies**: " + ", ".join(techs) + "\n"
        
        # Add search methodology
        formatted_response += f"""
## Search Methodology
**Search Type Used**: {search_mode.capitalize()} Search
**Extracted Terms**: {extracted_terms}
**Results Count**: {len(results)}

---
Search conducted by Project Discovery Specialist
"""
        
        return formatted_response
