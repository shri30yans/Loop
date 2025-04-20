"use client";
import { useState, useEffect } from "react";
import { Card } from "@nextui-org/card";
import { Button } from "@nextui-org/button";
import { Textarea } from "@nextui-org/input";
import { useAuthStore } from "../../lib/auth/authStore";
import { sendMessage, fetchChatHistory, type ChatMessage } from "./actions";

export default function ChatPage() {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [currentMessage, setCurrentMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const access_token = useAuthStore((state) => state.access_token);

  useEffect(() => {
    if (access_token) {
      loadChatHistory();
    }
  }, [access_token]);

  const loadChatHistory = async () => {
    try {
      const history = await fetchChatHistory(access_token!);
      setMessages(history);
    } catch (error) {
      console.error("Failed to load chat history:", error);
    }
  };

  const handleSendMessage = async () => {
    if (!currentMessage.trim() || !access_token) return;

    setIsLoading(true);
    
    try {
      const newMessage: ChatMessage = {
        id: Date.now().toString(),
        content: currentMessage,
        type: 'user',
        timestamp: new Date()
      };

      setMessages(prev => [...prev, newMessage]);
      setCurrentMessage("");

      const response = await sendMessage(currentMessage, access_token);
      setMessages(prev => [...prev, response]);
    } catch (error) {
      console.error("Failed to send message:", error);
    } finally {
      setIsLoading(false);
    }
  };

  if (!access_token) {
    return (
      <div className="h-[calc(100vh-4rem)] flex items-center justify-center">
        <Card className="p-4">
          <p className="text-lg">Please login to access the chat.</p>
        </Card>
      </div>
    );
  }

  return (
    <div className="h-[calc(100vh-4rem)] flex flex-col">
      {/* Messages Container */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((message) => (
          <div
            key={message.id}
            className={`flex ${
              message.type === 'user' ? 'justify-end' : 'justify-start'
            }`}
          >
            <div
              className={`max-w-[70%] px-4 py-2 rounded-lg ${
                message.type === 'user'
                  ? 'bg-gray-200 text-black'
                  : 'text-white'
              }`}
            >
              {message.content}
            </div>
          </div>
        ))}
      </div>

      {/* Input Container */}
      <Card className="m-4 p-4">
        <div className="flex gap-2">
          <Textarea
            value={currentMessage}
            onChange={(e) => setCurrentMessage(e.target.value)}
            placeholder="Type your message here..."
            minRows={1}
            maxRows={4}
            className="flex-1"
            onKeyDown={(e) => {
              if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                handleSendMessage();
              }
            }}
            isDisabled={isLoading}
          />
          <Button
            color="primary"
            isDisabled={!currentMessage.trim() || isLoading}
            onPress={handleSendMessage}
            isLoading={isLoading}
          >
            {isLoading ? "Sending..." : "Send"}
          </Button>
        </div>
      </Card>
    </div>
  );
}