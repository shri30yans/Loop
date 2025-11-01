"use client";
import { useState, useEffect, useRef, Suspense } from "react";
import { Card, CardBody } from "@nextui-org/card";
import { Button } from "@nextui-org/button";
import { Textarea } from "@nextui-org/input";
import { IoSendSharp } from "react-icons/io5";
import { useAuth, type AuthState } from "../../lib/auth/authStore";
import { 
  initializeSocket, 
  sendMessageWs, 
  clearConversation as clearConversationWs, 
  disconnectSocket,
  type ChatMessage,
  type ProjectResult,
  type WebSocketResponse,
  type WebSocketError
} from "./actions";

// Loading component for the chat interface
function ChatLoading() {
  return (
    <div className="flex flex-col items-center justify-center h-[600px] space-y-4">
      <div className="animate-spin text-2xl">⏳</div>
      <div>Loading chat interface...</div>
    </div>
  );
}

// Main chat content component
function ChatContent() {
  const [isConnected, setIsConnected] = useState(false);
  const [isTyping, setIsTyping] = useState(false);
  const [isReady, setIsReady] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [currentMessage, setCurrentMessage] = useState("");
  const [isSending, setIsSending] = useState(false);
  const [isReconnecting, setIsReconnecting] = useState(false);
  const access_token = useAuth((state: AuthState) => state.token);
  const reconnectAttempts = useRef(0);
  const maxReconnectAttempts = 3;

  // Set component as ready after mount and initialize messages
  useEffect(() => {
    setIsReady(true);
    if (access_token) {
      setMessages([
        {
          id: "1",
          content: "Hi! 👋 I'm Loop, your AI assistant. How can I help you today?",
          type: 'llm',
          timestamp: new Date().toISOString()
        },
        {
          id: "2",
          content: "Ask me to find any kind of projects.",
          type: 'llm',
          timestamp: new Date().toISOString()
        }
      ]);
    }
  }, [access_token]);

  // Initialize WebSocket connection only after component is ready and authenticated
  useEffect(() => {
    if (!isReady || !access_token) return;
    
    const connectWebSocket = () => {
      console.log("Attempting to connect to WebSocket...");
      console.log("WebSocket URL:", process.env.NEXT_PUBLIC_SEARCH_AGENT_URL);
      setIsReconnecting(true);
      const socket = initializeSocket({
        onConnect: () => {
          console.log("Successfully connected to search agent at:", process.env.NEXT_PUBLIC_SEARCH_AGENT_URL);
          setIsConnected(true);
          setIsReconnecting(false);
          reconnectAttempts.current = 0;
        },
        onDisconnect: () => {
          console.log("Disconnected from search agent. Will attempt reconnect if under max attempts.");
          setIsConnected(false);
          
          // Try to reconnect if we haven't exceeded max attempts
          if (reconnectAttempts.current < maxReconnectAttempts) {
            reconnectAttempts.current += 1;
            setTimeout(connectWebSocket, 1000 * reconnectAttempts.current);
          }
        },
        onResponse: (response: WebSocketResponse) => {
          console.log("Received response:", response);
          setIsTyping(false);

          const newMessage: ChatMessage = {
            id: Date.now().toString(),
            content: response.message,
            type: 'llm',
            timestamp: new Date().toISOString(),
            results: response.results
          };

          setMessages(prev => [...prev, newMessage]);
        },
        onError: (error: WebSocketError) => {
          console.error("Search error:", error);
          setIsTyping(false);

          const errorMessage: ChatMessage = {
            id: Date.now().toString(),
            content: error.message || "An error occurred",
            type: 'llm',
            timestamp: new Date().toISOString()
          };

          setMessages(prev => [...prev, errorMessage]);
        },
        onConversationCleared: () => {
          setMessages([
            {
              id: "1",
              content: "Hi! 👋 I'm Loop, your AI assistant. How can I help you today?",
              type: 'llm',
              timestamp: new Date().toISOString()
            },
            {
              id: "2",
              content: "I can help you with programming, answer questions, debug issues, and much more!",
              type: 'llm',
              timestamp: new Date().toISOString()
            }
          ]);
        }
      }, access_token);
      
      return socket;
    };
    
    // Start the connection
    const socket = connectWebSocket();

    return () => {
      disconnectSocket();
    };
  }, [isReady, access_token]);

  const handleSendMessage = () => {
    if (!currentMessage.trim() || !isConnected || !access_token) return;

    setIsSending(true);
    
    try {
      // Create a new user message
      const newMessage: ChatMessage = {
        id: Date.now().toString(),
        content: currentMessage,
        type: 'user',
        timestamp: new Date().toISOString()
      };

      // Add to messages state
      setMessages(prev => [...prev, newMessage]);
      setCurrentMessage("");

      // Show typing indicator
      setIsTyping(true);

      // Send to WebSocket server
      try {
        console.log("Sending message via WebSocket...");
        sendMessageWs(currentMessage, access_token, access_token);
        console.log("Message sent successfully");
      } catch (error) {
        console.error("WebSocket send error:", error);
        setMessages(prev => [...prev, {
          id: Date.now().toString(),
          content: "Failed to send message. Please check your connection.",
          type: 'llm',
          timestamp: new Date().toISOString()
        }]);
      }
    } catch (error) {
      console.error("Failed to send message:", error);
      setIsTyping(false);
    } finally {
      setIsSending(false);
    }
  };

  const clearConversation = () => {
    if (!isConnected || !access_token) return;
    clearConversationWs(access_token, access_token);
  };

  // Scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  // Show loading state if not ready, and handle redirect if not authenticated
  if (!isReady) {
    return <ChatLoading />;
  }
  
  if (!access_token) {
    if (typeof window !== 'undefined') {
      window.location.href = '/auth/login';
    }
    return <ChatLoading />;
  }

  return (
    <div className="flex flex-col h-full relative max-w-3xl mx-auto">
      {/* Header with Clear Conversation */}
      <div className="absolute top-0 right-0 p-4">
        <Button
          size="sm"
          variant="light"
          onPress={clearConversation}
          isDisabled={!isConnected}
        >
          Clear Conversation
        </Button>
      </div>
      {/* Messages Container */}
      <div className="flex-1 space-y-8 overflow-y-auto absolute top-[48px] left-0 right-0 bottom-[100px] pr-2" key="messages-container">
        <div className="px-6 space-y-4 pb-8">
        {!isConnected && (
          <div className="text-center text-red-500 py-2">
            Disconnected from search service. Please try refreshing the page.
          </div>
        )}
        {messages.map((message) => (
          <div
            key={message.id}
            className={`flex ${
              message.type === 'user' ? 'justify-end' : 'justify-start'
            }`}
          >
            <div
              className={`max-w-[60%] ${
                message.type === 'user' ? 'ml-auto' : 'mr-auto'
              }`}
            >
              {message.type === 'user' ? (
                <Card
                  className="px-4 py-2 text-md bg-primary text-primary-foreground rounded-3xl shadow-md transition-all duration-200"
                >
                  {message.content}
                </Card>
              ) : (
                <div className="space-y-4">
                  <div className="px-2 text-md text-foreground">
                    {message.content}
                  </div>
                  {message.results && message.results.length > 0 && (
                    <div className="mt-4 space-y-4">
                      {message.results.map((result, i) => (
                        <Card key={i} className="p-4 bg-background/40 backdrop-blur-sm">
                          <h3 className="font-semibold text-lg">
                            {result.project.title}
                          </h3>
                          <p className="text-sm text-default-500 mt-1">
                            {result.project.summary}
                          </p>
                          <div className="mt-2 flex items-center">
                            <div className="text-xs text-default-400">
                              Relevance: {Math.round(result.score * 100)}%
                            </div>
                            <Button
                              as="a"
                              href={`/projects/${result.project.id}`}
                              className="ml-auto text-sm"
                              color="primary"
                              variant="light"
                              size="sm"
                            >
                              View Project →
                            </Button>
                          </div>
                        </Card>
                      ))}
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        ))}
        <div ref={messagesEndRef} />
        </div>
      </div>

      {/* Input Container */}
      <Card className=" absolute bottom-0 left-0 right-0 rounded-3xl bg-background/80 backdrop-blur-md" shadow="lg">
        <CardBody>
          <div className="flex gap-4 items-stretch">
          <Textarea
            value={currentMessage}
            onChange={(e) => setCurrentMessage(e.target.value)}
            placeholder="Type your message here..."
            minRows={2}
            maxRows={2}
            className="flex-1 text-xl rounded-xl "
            variant="faded"
            onKeyDown={(e) => {
              if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                handleSendMessage();
              }
            }}
            isDisabled={isSending}
          />
          <Button
            color="primary"
            size="lg"
            className=" rounded-xl hover:scale-105 hover:shadow-lg active:scale-95 transition-all duration-200 flex items-center justify-center bg-gradient-to-tr from-primary-500 to-primary-600"
            isDisabled={!currentMessage.trim() || isSending}
            onPress={handleSendMessage}
            isLoading={isSending}
          >
            {isSending ? (
              <div className="animate-spin">⏳</div>
            ) : (
              <IoSendSharp size={20} />
            )}
          </Button>
          </div>
        </CardBody>
      </Card>
    </div>
  );
}

// Export the page component with proper loading state
export default function ChatPage() {
  return (
    <Suspense fallback={<ChatLoading />}>
      <ChatContent />
    </Suspense>
  );
}
