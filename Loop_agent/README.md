# Loop Agent Service

The Loop Agent Service is an intelligent project discovery system that uses a team of specialized agents to help users find relevant software projects.

## Architecture

The service uses an agent-based architecture with the following components:

### Agents

1. **Project Search Agent**
   - Extracts key search terms from user queries
   - Determines optimal search strategy (vector/graph/hybrid)
   - Coordinates with Go backend for search execution

2. **Clarification Agent**
   - Analyzes query ambiguity
   - Generates contextual clarifying questions
   - Helps refine search parameters

3. **Deep Search Agent**
   - Performs comprehensive analysis of search results
   - Identifies patterns and relationships
   - Provides technical insights and recommendations

### Team Coordination

The `ProjectDiscoveryTeam` class orchestrates the interactions between agents:
- Routes queries to appropriate agents
- Manages conversation history
- Coordinates deep search analysis
- Formats and presents results

## Installation

### For Development

1. Clone the repository:
```bash
git clone <repository-url>
cd Loop_agent
```

2. Create and activate a virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

3. Install in development mode:
```bash
pip install -e .
```

4. Copy and configure environment variables:
```bash
cp .env.example .env
# Edit .env with your settings
```

### For Production

1. Install from PyPI:
```bash
pip install loop-agent
```

2. Configure environment variables as needed.

## Usage

### Running the Service

1. Start the service:
```bash
# Using the module
python -m loop_agent

# Or using the app directly
python app.py
```

2. The service will be available at `http://localhost:5000`

### API Endpoints

#### HTTP Endpoints

- `GET /health`: Health check endpoint
- `POST /embed`: Generate embedding for a single text
- `POST /embed-batch`: Generate embeddings for multiple texts

#### WebSocket Events

- `connect`: Client connection event
- `disconnect`: Client disconnection event
- `search`: Handle search requests
  ```javascript
  // Example search request
  socket.emit('search', {
    query: "Find React projects using GraphQL",
    userId: "user123",
    deepSearch: true
  });
  ```
- `clear_conversation`: Clear conversation history
  ```javascript
  socket.emit('clear_conversation', {
    userId: "user123"
  });
  ```

## Development

### Project Structure

```
Loop_agent/
├── __init__.py           # Package initialization
├── __main__.py          # Entry point for running as module
├── app.py               # Main Flask application
├── requirements.txt     # Package dependencies
├── setup.py            # Package setup configuration
├── services/           # Core services
│   ├── __init__.py
│   ├── embedding.py    # Embedding service
│   ├── go_client.py    # Go backend client
│   └── agents/         # Agent implementations
│       ├── __init__.py
│       ├── base.py     # Base agent class
│       ├── project_search.py
│       ├── clarification.py
│       └── deep_search.py
└── team/              # Team coordination
    ├── __init__.py
    └── team.py       # Team implementation
```

### Adding New Agents

1. Create a new agent class in `services/agents/`
2. Inherit from `BaseAgent`
3. Implement required methods
4. Update `ProjectDiscoveryTeam` to use the new agent

### Extending Search Capabilities

1. Add new search modes in `ProjectSearchAgent`
2. Update the Go backend client if needed
3. Modify the team coordination logic

## Configuration

### Environment Variables

```bash
# API Configuration
API_KEY=your_groq_api_key
API_URL=https://api.groq.com/openai
MODEL_NAME=meta-llama/llama-4-maverick-17b-128e-instruct

# Embedding Model
EMBEDDING_MODEL=all-MiniLM-L6-v2

# Backend Configuration
GO_BACKEND_URL=http://localhost:8080

# Server Configuration
PORT=5000
HOST=0.0.0.0

# WebSocket Configuration
SOCKET_PING_TIMEOUT=60
SOCKET_PING_INTERVAL=25
```

## ASGI Implementation

The Loop Agent uses an ASGI (Asynchronous Server Gateway Interface) architecture for high-performance WebSocket communication. This implementation provides several key benefits:

### Performance Benefits

- **High Concurrency**: Efficiently handles thousands of concurrent WebSocket connections
- **Reduced Latency**: Asynchronous processing minimizes response times for I/O-bound operations
- **Resource Efficiency**: Non-blocking I/O maximizes CPU and memory utilization
- **Scalability**: Scales horizontally across multiple workers

### Technical Implementation

- Uses `flask-socketio` with `async_mode='asgi'`
- Deployed with Uvicorn, an ASGI server optimized for WebSockets
- Implements proper async/await patterns throughout the codebase
- Utilizes connection pooling for backend service communication

### Running in Production

For production deployments, use the provided Dockerfile which configures Uvicorn with optimized settings:

```bash
docker build -t loop-agent .
docker run -p 5000:5000 loop-agent
```

The production configuration uses:
- Multiple worker processes based on available CPU cores
- Optimized event loop with uvloop
- Fast HTTP parsing with httptools
- Native WebSocket implementation

### Performance Tuning

For high-traffic deployments, consider the following configuration options:

1. **Worker Count**: Adjust the number of Uvicorn workers in the Dockerfile based on available CPU cores
   ```
   # For a machine with 4 cores, use 8 workers
   CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "5000", "--workers", "8", ...]
   ```

2. **Connection Limits**: Adjust the connection pool limits in `go_client.py` based on expected load
   ```python
   connector = aiohttp.TCPConnector(limit=200)  # Increase for higher concurrency
   ```

3. **WebSocket Timeouts**: Configure ping intervals and timeouts in `.env`:
   ```
   SOCKET_PING_TIMEOUT=120
   SOCKET_PING_INTERVAL=30
   ```

## Error Handling

The service includes comprehensive error handling:
- Invalid queries
- Backend communication errors
- LLM API errors
- Search execution errors

Each error is logged and returned with appropriate context.

## Logging

The service uses Python's logging module with:
- INFO level for normal operations
- ERROR level for issues
- Detailed context for debugging

Logs include:
- Query processing steps
- Agent interactions
- Search execution details
- Error traces

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

Please include tests and documentation updates.

## License

[Add your license information here]
