import { io, Socket } from "socket.io-client";

export interface ProjectResult {
  project: {
    id: string;
    title: string;
    summary: string;
  };
  score: number;
  metadata?: any;
}

export interface ChatMessage {
  id: string;
  content: string;
  type: "user" | "llm";
  timestamp: string;
  results?: ProjectResult[];
}

export interface WebSocketResponse {
  type: string;
  message: string;
  results?: ProjectResult[];
  partial_results?: ProjectResult[];
}

export interface WebSocketError {
  message: string;
}

export interface ChatSocketEvents {
  onConnect: () => void;
  onDisconnect: () => void;
  onResponse: (response: WebSocketResponse) => void;
  onError: (error: WebSocketError) => void;
  onConversationCleared: () => void;
}

let socket: Socket | null = null;

/**
 * Initialize WebSocket connection to the chat server
 */
export function initializeSocket(
  events: ChatSocketEvents,
  access_token?: string
): Socket {
  if (socket) {
    return socket;
  }

  // Create new socket connection with explicit configuration
  socket = io(process.env.NEXT_PUBLIC_SEARCH_AGENT_URL || 'http://localhost:5000', {
    transports: ['websocket'],  // Force WebSocket only, no polling
    reconnectionAttempts: 5,
    reconnectionDelay: 1000,
    autoConnect: true,
    forceNew: true,  // Force a new connection
    path: '/socket.io',  // Default Socket.IO path
    withCredentials: false,  // Don't send credentials
    extraHeaders: {
      'Access-Control-Allow-Origin': '*'
    },
    // Include auth token if available
    auth: access_token ? {
      token: access_token
    } : undefined,
    // Debug options
    timeout: 10000,  // Increase timeout
    rememberUpgrade: true,
    rejectUnauthorized: false
  });

  // Debug logging for connection issues
  socket.io.on("error", (error) => {
    console.error("Socket.IO manager error:", error);
  });

  socket.io.on("reconnect_attempt", (attempt) => {
    console.log("Socket.IO reconnection attempt:", attempt);
  });

  socket.io.on("reconnect_error", (error) => {
    console.error("Socket.IO reconnection error:", error);
  });

  // Monitor connection health
  socket.io.engine?.on("ping", () => {
    console.log("Socket.IO ping sent");
  });

  socket.io.engine?.on("packet", (packet) => {
    if (packet.type === 'pong') {
      console.log("Socket.IO pong received");
    }
  });

  // Debug logging
  socket.on('connect_error', (error) => {
    console.error('Socket.IO connection error:', error);
  });

  socket.on('connect_timeout', () => {
    console.error('Socket.IO connection timeout');
  });

  // Set up event listeners
  socket.on("connect", events.onConnect);
  socket.on("disconnect", events.onDisconnect);
  socket.on("response", events.onResponse);
  socket.on("error", events.onError);
  socket.on("conversation_cleared", events.onConversationCleared);

  return socket;
}

/**
 * Send a message to the chat server via WebSocket
 */
export function sendMessageWs(
  message: string,
  userId: string,
  access_token?: string
): void {
  if (!socket || !socket.connected) {
    throw new Error("WebSocket not connected");
  }

  socket.emit("search", {
    query: message,
    userId: userId,
    token: access_token // Include token in the payload
  });
}

/**
 * Clear conversation history
 */
export function clearConversation(userId: string, access_token?: string): void {
  if (!socket || !socket.connected) {
    throw new Error("WebSocket not connected");
  }

  socket.emit("clear_conversation", {
    userId: userId,
    token: access_token // Include token in the payload
  });
}

/**
 * Disconnect WebSocket
 */
export function disconnectSocket(): void {
  if (socket) {
    socket.disconnect();
    socket = null;
  }
}
