# Loop Agent WebSocket API Documentation

This document outlines the WebSocket event schemas for communication between the frontend and the Loop agent service.

## Connection

The WebSocket server is available at: `ws://localhost:5000`

## Performance Considerations

The Loop Agent uses an ASGI (Asynchronous Server Gateway Interface) setup with SocketIO for high-performance WebSocket communication. This architecture provides several benefits:

1. **Improved Concurrency**: ASGI allows handling thousands of concurrent WebSocket connections efficiently
2. **Reduced Latency**: Asynchronous processing minimizes response times for I/O-bound operations
3. **Better Resource Utilization**: Non-blocking I/O maximizes CPU and memory efficiency
4. **Scalability**: The architecture scales well horizontally across multiple workers

For production deployments, the service uses Uvicorn with optimized settings for WebSocket performance.

## Authentication

Authentication is handled via tokens passed during connection or in specific events.

### Connection with Authentication

```
ws://localhost:5000?auth=<token>
```

Alternatively, you can pass the token in the Authorization header:

```
Authorization: Bearer <token>
```

## Events

### 1. Connect

**Event**: Automatic on WebSocket connection

**Client Payload**: None (Authentication can be provided in the connection URL or headers)

**Server Response**:
```json
{
  "status": "connected",
  "sid": "socket-session-id"
}
```

**Error Response**:
```json
{
  "status": "error",
  "message": "Connection failed",
  "error": "Error details"
}
```

### 2. Register

**Event**: `register`

**Description**: Registers a user's socket session with the agent service

**Client Payload**:
```json
{
  "userId": "user-uuid",
  "token": "jwt-token" // Optional if provided during connection
}
```

**Server Response**:
```json
{
  "status": "success"
}
```

**Error Response**:
```json
{
  "message": "Failed to register session",
  "error": "Error details"
}
```

### 3. Search

**Event**: `search`

**Description**: Performs a search query using the agent service

**Client Payload**:
```json
{
  "query": "search query string",
  "userId": "user-uuid",
  "deepSearch": false, // Optional, defaults to false
  "token": "jwt-token" // Optional if provided during connection
}
```

**Server Response** (Event: `response`):
```json
{
  "type": "results",
  "message": "Formatted search results with analysis",
  "results": [
    {
      "name": "Project Name",
      "description": "Project description",
      "score": 0.95,
      "technologies": ["React", "Node.js"],
      "url": "project-url",
      "stars": 123,
      "last_updated": "2023-01-01",
      "language": "JavaScript"
    }
  ],
  "search_mode": "vector|graph|hybrid"
}
```

**Clarification Response** (When more information is needed):
```json
{
  "type": "clarification",
  "message": "To help you find the most relevant projects, I need some clarification:\n\n- Question 1\n- Question 2",
  "partial_results": [],
  "analysis": "Analysis of the query"
}
```

**Non-Search Response** (For conversational queries):
```json
{
  "type": "non_search",
  "message": "Conversational response",
  "results": []
}
```

**Error Response**:
```json
{
  "message": "An error occurred while processing your request",
  "error": "Error details"
}
```

### 4. Clear Conversation

**Event**: `clear_conversation`

**Description**: Clears the conversation history for a user

**Client Payload**:
```json
{
  "userId": "user-uuid",
  "token": "jwt-token" // Optional if provided during connection
}
```

**Server Response**:
```json
{
  "status": "success"
}
```

**Error Response**:
```json
{
  "message": "Failed to clear conversation",
  "error": "Error details"
}
```

### 5. Disconnect

**Event**: Automatic on WebSocket disconnection

**Client Payload**: None

**Server Response**: None

## Error Handling

All events can return an error response with the following structure:

```json
{
  "message": "Human-readable error message",
  "error": "Technical error details"
}
```

## Search Modes

The agent supports three search modes:

1. **Vector Search**: Finds semantically similar projects based on natural language understanding
2. **Graph Search**: Finds projects based on specific attributes and relationships
3. **Hybrid Search**: Combines vector and graph search for complex queries

The search mode is automatically determined based on the query, but can be influenced by the query structure:

- "Show me projects similar to Netflix" → Vector search
- "Find React projects using GraphQL" → Graph search
- "Projects like Airbnb but using Vue.js" → Hybrid search

## Response Types

The search response can be one of three types:

1. **results**: Standard search results
2. **clarification**: When the agent needs more information to provide better results
3. **non_search**: For conversational queries that don't require search

## Example Usage

### Connecting and Registering

```javascript
const socket = io('http://localhost:5000', {
  transports: ['websocket'],  // Prefer WebSocket transport for better performance
  upgrade: false,             // Disable transport upgrades for consistent performance
  timeout: 60000              // Increase timeout for long-running operations
});

socket.on('connect', () => {
  console.log('Connected to agent service');
  
  // Register the user
  socket.emit('register', {
    userId: 'user-123',
    token: 'jwt-token'
  });
});

socket.on('registered', (data) => {
  console.log('Registration successful', data);
});
```

### Performing a Search

```javascript
// Send a search query
socket.emit('search', {
  query: 'Find React projects with GraphQL',
  userId: 'user-123'
});

// Listen for the response
socket.on('response', (data) => {
  if (data.type === 'results') {
    console.log('Search results:', data.results);
  } else if (data.type === 'clarification') {
    console.log('Clarification needed:', data.message);
  } else {
    console.log('Conversational response:', data.message);
  }
});
```

### Clearing Conversation History

```javascript
socket.emit('clear_conversation', {
  userId: 'user-123'
});

socket.on('conversation_cleared', (data) => {
  console.log('Conversation cleared:', data);
});
```

### Error Handling

```javascript
socket.on('error', (error) => {
  console.error('Error from agent service:', error);
});
```

## Scaling Considerations

When deploying the Loop Agent in production environments, consider the following scaling strategies:

1. **Worker Configuration**: Adjust the number of Uvicorn workers based on available CPU cores (typically 2x number of cores)

2. **Connection Pooling**: The agent uses connection pooling for backend services to minimize resource usage

3. **Load Balancing**: For high-traffic deployments, use a load balancer in front of multiple agent instances

4. **Memory Management**: Monitor memory usage as WebSocket connections can consume significant resources at scale

5. **Timeouts and Heartbeats**: The service uses configurable ping intervals and timeouts to maintain connection health

For very high concurrency scenarios (10,000+ concurrent connections), consider deploying the service on machines with:
- At least 4 CPU cores
- 8GB+ RAM
- Sufficient network bandwidth
- Properly configured OS limits for open files and network connections