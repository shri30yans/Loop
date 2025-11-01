"""
Main Flask application for the Loop agent service.
"""

import asyncio
from flask import Flask, request, jsonify
from flask_cors import CORS
from flask_socketio import SocketIO, emit
from services.embedding import EmbeddingService
from services.go_client import GoBackendClient
from team import ProjectDiscoveryTeam
import logging
import time
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Initialize Flask app
app = Flask(__name__)
CORS(app)

# Configure Socket.IO
socketio = SocketIO(
    app,
    cors_allowed_origins="*",
    async_mode='asgi',  # Using ASGI mode for better performance with WebSockets
    ping_timeout=int(os.getenv('SOCKET_PING_TIMEOUT', '60')),
    ping_interval=int(os.getenv('SOCKET_PING_INTERVAL', '25')),
    logger=True,
    engineio_logger=True
)

# Initialize services
model_name = os.getenv('EMBEDDING_MODEL', 'all-MiniLM-L6-v2')
go_backend_url = os.getenv('GO_BACKEND_URL', 'http://localhost:8080')

embedding_service = EmbeddingService(model_name=model_name)
go_client = GoBackendClient(go_backend_url)
discovery_team = ProjectDiscoveryTeam(go_client)

# Dictionary to store tokens by session ID
session_tokens = {}

@app.route('/health', methods=['GET'])
def health_check():
    """Health check endpoint."""
    return jsonify({
        'status': 'healthy',
        'model': model_name,
        'dimensions': embedding_service.get_dimensions(),
        'go_backend': go_backend_url
    })

@app.route('/embed', methods=['POST'])
def generate_embedding():
    """Generate embedding for a single text."""
    start_time = time.time()
    
    # Get request data
    data = request.json
    if not data or 'text' not in data:
        return jsonify({'error': 'Missing required field: text'}), 400
    
    text = data['text']
    
    try:
        # Generate embedding
        embedding = embedding_service.generate_embedding(text)
        
        # Convert to list for JSON serialization
        embedding_list = embedding.tolist()
        
        # Calculate processing time
        processing_time = time.time() - start_time
        
        return jsonify({
            'embedding': embedding_list,
            'dimensions': len(embedding_list),
            'processing_time_ms': round(processing_time * 1000, 2)
        })
    except Exception as e:
        logger.error(f"Error generating embedding: {str(e)}")
        return jsonify({'error': str(e)}), 500

@app.route('/embed-batch', methods=['POST'])
def generate_batch_embeddings():
    """Generate embeddings for multiple texts."""
    start_time = time.time()
    
    # Get request data
    data = request.json
    if not data or 'texts' not in data:
        return jsonify({'error': 'Missing required field: texts'}), 400
    
    texts = data['texts']
    
    try:
        # Generate embeddings
        embeddings = embedding_service.generate_batch_embeddings(texts)
        
        # Convert to list for JSON serialization
        embeddings_list = [emb.tolist() for emb in embeddings]
        
        # Calculate processing time
        processing_time = time.time() - start_time
        
        return jsonify({
            'embeddings': embeddings_list,
            'count': len(embeddings_list),
            'dimensions': embedding_service.get_dimensions(),
            'processing_time_ms': round(processing_time * 1000, 2)
        })
    except Exception as e:
        logger.error(f"Error generating batch embeddings: {str(e)}")
        return jsonify({'error': str(e)}), 500

# WebSocket event handlers
@socketio.on('connect')
async def handle_connect():
    """Handle client connection."""
    logger.info(f"Client connected from {request.sid}")
    
    # Extract token from auth data
    auth = request.args.get('auth') or {}
    if isinstance(auth, dict):
        token = auth.get('token')
    else:
        # Try to get from headers
        token = request.headers.get('Authorization')
        if token and token.startswith('Bearer '):
            token = token.split(' ')[1]
    
    # Store token if available
    if token:
        logger.info(f"Received auth token for session {request.sid}")
        session_tokens[request.sid] = token
    else:
        logger.warning(f"No auth token provided for session {request.sid}")
    
    await emit('connect_response', {'status': 'connected', 'sid': request.sid})
    return True

@socketio.on('register')
async def handle_register(data):
    """Register a user's socket session."""
    try:
        user_id = data.get('userId', 'anonymous')
        token = data.get('token') or session_tokens.get(request.sid)
        
        discovery_team.register_session(request.sid, user_id, token)
        await socketio.emit('registered', {'status': 'success'}, room=request.sid)
    except Exception as e:
        logger.error(f"Error registering session: {str(e)}")
        await socketio.emit('error', {
            'message': 'Failed to register session',
            'error': str(e)
        }, room=request.sid)

@socketio.on('disconnect')
async def handle_disconnect():
    """Handle client disconnection."""
    logger.info("Client disconnected")

@socketio.on('search')
async def handle_search(data):
    """Handle search requests using the agent team."""
    try:
        query = data.get('query', '')
        user_id = data.get('userId', 'anonymous')
        deep_search = data.get('deepSearch', False)
        token = data.get('token') or session_tokens.get(request.sid)
        
        # Capture the session ID
        sid = request.sid
        
        # Ensure the session is registered
        if sid not in discovery_team.sid_to_user:
            discovery_team.register_session(sid, user_id, token)
        
        logger.info(f"Search request from {user_id}: {query}")
        
        try:
            # Process query and get response asynchronously
            response = await discovery_team.process_query(query, user_id, deep_search, token)
            # Send response back to client
            await socketio.emit('response', response, room=sid)
        except Exception as e:
            logger.error(f"Error processing query: {str(e)}")
            await socketio.emit('error', {
                'message': 'An error occurred while processing your request',
                'error': str(e)
            }, room=sid)
        finally:
            await go_client.close()
            
    except Exception as e:
        logger.error(f"Error in search handler: {str(e)}")
        await socketio.emit('error', {
            'message': 'An error occurred while processing your request',
            'error': str(e)
        }, room=request.sid)

@socketio.on('clear_conversation')
async def handle_clear(data):
    """Clear conversation history for a user."""
    try:
        user_id = data.get('userId', 'anonymous')
        token = data.get('token') or session_tokens.get(request.sid)
        sid = request.sid
        
        # Ensure the session is registered
        if sid not in discovery_team.sid_to_user:
            discovery_team.register_session(sid, user_id, token)
            
        discovery_team.clear_conversation(user_id)
        await socketio.emit('conversation_cleared', {'status': 'success'}, room=sid)
    except Exception as e:
        logger.error(f"Error clearing conversation: {str(e)}")
        await socketio.emit('error', {
            'message': 'Failed to clear conversation',
            'error': str(e)
        }, room=request.sid)

if __name__ == '__main__':
    # Get port from environment variable or use default
    port = int(os.getenv('PORT', 5000))
    host = os.getenv('HOST', '0.0.0.0')
    
    # Run the Flask app with SocketIO using ASGI
    logger.info(f"Starting service with ASGI WebSocket support on {host}:{port}")
    socketio.run(app, host=host, port=port)
    # Note: In production, use an ASGI server like Uvicorn or Hypercorn
    # The socketio.run is only for development
