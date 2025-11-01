# Loop Technology Stack

## Backend Services

### Main Backend (Go)
- **Framework**: Standard Go HTTP with Gorilla Mux
- **Architecture**: Clean Architecture with layers:
  - Handlers (API endpoints)
  - Services (Business logic)
  - Repositories (Data access)
  - Models (Domain entities)

### Databases
- **PostgreSQL**
  - Primary relational database
  - pgvector extension for vector similarity search
  - Stores user data, projects, and vector embeddings
- **Neo4j**
  - Graph database for knowledge graph
  - Stores entities and relationships
  - Enables graph-based search

### Embedding Service (Python)
- **Framework**: Flask
- **Core Libraries**:
  - sentence-transformers for text embeddings
  - numpy for numerical operations
- **Features**:
  - Single text embedding generation
  - Batch embedding processing
  - Health monitoring
  - Model configuration options

## Frontend

### Next.js Application
- **Framework**: Next.js 13+
- **UI Components**: NextUI
- **Styling**: Tailwind CSS
- **State Management**: React Context
- **Type Safety**: TypeScript

## AI and NLP Components

### AI Integration
- **Providers**:
  - OpenAI API
  - Ollama (local models)
- **Features**:
  - Entity extraction
  - Summary generation
  - Query transformation

### Search System
- **Hybrid Search**:
  - Vector similarity search (pgvector)
  - Graph-based search (Neo4j)
  - Configurable weight system
- **Embedding Models**:
  - sentence-transformers (default: all-MiniLM-L6-v2)
  - Support for alternative models

## Infrastructure

### Development
- **Version Control**: Git
- **Package Management**:
  - Go Modules
  - npm/yarn (Frontend)
  - pip (Python)
- **Development Tools**:
  - VSCode
  - Docker (optional)
  - Make (build automation)

### API Design
- **Protocol**: REST
- **Authentication**: JWT
- **Documentation**: OpenAPI/Swagger

## Key Technical Decisions

### Vector Search Implementation
- Chose pgvector over dedicated vector databases for:
  - Simplified infrastructure
  - Integration with existing PostgreSQL
  - Sufficient performance for current scale
  - ACID compliance

### Embedding Service Architecture
- Separate Python service for:
  - Better ML library compatibility
  - Independent scaling
  - Flexibility in model selection
  - Easy local development

### Hybrid Search Approach
- Combined graph and vector search for:
  - Better semantic understanding
  - Context-aware results
  - Configurable search behavior
  - Improved result relevance

## Future Considerations

### Scalability
- Potential migration to dedicated vector database
- Distributed embedding service
- Caching layer for search results

### Performance
- Search result caching
- Batch processing optimizations
- Index optimization

### Features
- Real-time search suggestions
- Advanced query understanding
- Personalized search rankings
- Analytics and monitoring
