# Loop - Project Discovery Platform

Loop is a platform where builders share their projects and the journey of building them. It enables users to discover projects through natural language queries, leveraging NLP and knowledge graph technologies.

## Features

- Natural language project search with AI-powered understanding
- Interactive chat interface for project discovery
- Knowledge graph-based project relationships
- Vector similarity search for semantic matching
- Real-time search suggestions and clarifications
- User authentication and profiles
- Project creation and management

## Architecture

The platform consists of three main components:

1. **Frontend (Next.js)**
   - Modern React-based UI with TypeScript
   - Real-time WebSocket communication
   - Responsive chat interface
   - Authentication and user management

2. **Backend (Go)**
   - RESTful API endpoints
   - PostgreSQL for relational data
   - Neo4j for knowledge graph
   - JWT-based authentication
   - Project management and analysis

3. **Embedding Service (Python)**
   - Flask-based WebSocket server
   - Text embedding generation
   - Vector similarity search
   - Interactive search agent
   - Real-time chat handling

## Prerequisites

- Go 1.19 or later
- Node.js 16 or later
- Python 3.8 or later
- PostgreSQL 14 or later with pgvector extension
- Neo4j 4.4 or later

## Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/loop.git
   cd loop
   ```

2. **Install dependencies**

   Frontend:
   ```bash
   cd Loop_frontend
   npm install
   cd ..
   ```

   Embedding Service:
   ```bash
   cd embedding_service
   pip install -r requirements.txt
   cd ..
   ```

   Backend:
   ```bash
   cd Loop_backend
   go mod download
   cd ..
   ```

3. **Configure environment variables**

   Create `.env.local` in Loop_frontend:
   ```
   NEXT_PUBLIC_SEARCH_AGENT_URL=http://localhost:5000
   NEXT_PUBLIC_BASE_API_URL=http://localhost:8080
   ```

   Create `.env` in embedding_service:
   ```
   PORT=5000
   EMBEDDING_MODEL=all-MiniLM-L6-v2
   GO_BACKEND_URL=http://localhost:8080
   CORS_ORIGINS=http://localhost:3000
   ```

4. **Set up databases**

   PostgreSQL:
   ```sql
   CREATE EXTENSION vector;
   ```

   Run migrations:
   ```bash
   cd Loop_backend
   go run migrations/main.go
   ```

## Running the Project

You can run all services together using the development script:

```bash
./dev.sh
```

Or run each service individually:

1. **Start the Go backend**
   ```bash
   cd Loop_backend
   go run main.go
   ```

2. **Start the embedding service**
   ```bash
   cd embedding_service
   python app.py
   ```

3. **Start the frontend**
   ```bash
   cd Loop_frontend
   npm run dev
   ```

The application will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Embedding Service: http://localhost:5000

## Development

- Frontend code is in `Loop_frontend/`
- Backend code is in `Loop_backend/`
- Embedding service code is in `embedding_service/`

## Testing

Run tests for each component:

```bash
# Backend tests
cd Loop_backend
go test ./...

# Frontend tests
cd Loop_frontend
npm test

# Embedding service tests
cd embedding_service
python -m pytest
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
