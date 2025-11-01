# Loop Codebase Summary

## Project Structure Overview

The Loop project is organized into three main components:

1. **Loop_frontend**: Next.js frontend application
2. **Loop_backend**: Go backend service
3. **embedding_service**: Python Flask service for generating embeddings
4. **graphrag**: Additional NLP-related utilities

## Key Components and Their Interactions

### Backend Components

#### Core Services
- **AuthService**: Handles user authentication and JWT management
- **UserService**: Manages user profiles and data
- **ProjectService**: Handles project creation, retrieval, and updates
- **QueryService**: Transforms natural language queries to Cypher queries
- **SummaryService**: Generates project summaries and manages vector storage
- **EmbeddingService**: Interfaces with the Python embedding service
- **HybridSearchService**: Combines vector and graph search results with configurable weights

#### AI Integration
- **Provider Interface**: Abstract interface for AI providers
- **Provider Implementations**: 
  - OpenAI provider
  - Ollama provider (local model support)
- **ProjectProcessor**: Analyzes projects using AI to extract entities and summaries
- **EmbeddingClient**: Communicates with the Python embedding service

#### Data Access
- **Repositories**: 
  - AuthRepository: Authentication data
  - UserRepository: User data
  - ProjectRepository: Project data
  - GraphRepository: Knowledge graph operations
  - VectorRepository: Vector storage and similarity search with pgvector

#### API Handlers
- **AuthHandler**: Authentication endpoints
- **UserHandler**: User management endpoints
- **ProjectHandler**: Project CRUD operations
- **SearchHandler**: Hybrid search functionality with configurable weights
- **SummaryHandler**: Project summary and vector storage endpoints

### Embedding Service Components
- **Flask Application**: Lightweight web server
- **Embedding Service**: Generates embeddings using sentence-transformers
- **API Endpoints**: 
  - `/embed`: Single text embedding generation
  - `/embed-batch`: Batch embedding generation
  - `/health`: Service health check

### Frontend Components

#### Pages
- Authentication pages (login, register)
- Project pages (creation, viewing, listing)
- Search interface with weight adjustment
- User profile pages

#### UI Components
- Navigation components
- Project cards
- Advanced search interface with weight controls
- Forms for project creation and updates

## Data Flow

### Project Creation Flow
1. User submits project data via frontend form
2. ProjectHandler receives request and validates data
3. ProjectService creates project in PostgreSQL
4. ProjectProcessor analyzes project content using AI
5. Extracted entities are stored in Neo4j as a knowledge graph
6. Project summary is generated and sent to embedding service
7. Embedding service returns vector representation
8. Vector, summary, and extracted metadata are stored in project_vectors table

### Hybrid Search Flow
1. User enters natural language query with optional weight adjustments
2. Query is sent to SearchHandler
3. Query is embedded via embedding service
4. VectorRepository performs similarity search using pgvector
5. GraphService performs entity-based search in Neo4j
6. Results are combined using configurable weights:
   - Vector similarity score * vectorWeight
   - Graph relevance score * graphWeight
7. Final results are ranked and returned to frontend

## External Dependencies

### Backend Dependencies
- **github.com/gorilla/mux**: HTTP routing
- **github.com/joho/godotenv**: Environment variable management
- **github.com/rs/cors**: CORS handling
- **Neo4j driver**: Graph database connectivity
- **PostgreSQL driver**: Relational database connectivity
- **pgvector extension**: Vector similarity search in PostgreSQL

### Embedding Service Dependencies
- **Flask**: Python web framework
- **sentence-transformers**: Embedding generation
- **numpy**: Numerical operations
- **Docker**: Containerization (optional)

### Frontend Dependencies
- **Next.js**: React framework
- **NextUI**: UI component library
- **Tailwind CSS**: Utility-first CSS framework
- **TypeScript**: Type-safe JavaScript

## Recent Significant Changes
- Implementation of hybrid search system combining graph and vector search
- Addition of configurable weights for search components
- Creation of project_vectors table with pgvector support
- Development of Python embedding service
- Implementation of project summary generation and storage
- Migration script for existing projects
- Fixed vector format issue in PostgreSQL storage (proper formatting for pgvector)

## User Feedback Integration
- No documented user feedback integration yet (early development stage)
