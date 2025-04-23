"use client";
import { useState, useEffect, useRef } from "react";
import { Card, CardHeader, CardBody } from "@nextui-org/card";
import { Button } from "@nextui-org/button";
import { Textarea } from "@nextui-org/input";
import { IoSendSharp } from "react-icons/io5";
import { useAuthStore } from "../../lib/auth/authStore";
import { sendMessage, fetchChatHistory, type ChatMessage } from "./actions";

export default function ChatPage() {
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const initialMessages: ChatMessage[] = [
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
  ];

  const [messages, setMessages] = useState<ChatMessage[]>(initialMessages);
  const [currentMessage, setCurrentMessage] = useState("");
  const [isPageLoading, setIsPageLoading] = useState(true);
  const [isSending, setIsSending] = useState(false);
  const access_token = useAuthStore((state) => state.access_token);

  useEffect(() => {
    if (access_token) {
      loadChatHistory().finally(() => setIsPageLoading(false));
    } else {
      setIsPageLoading(false);
    }
  }, [access_token]);

  const loadChatHistory = async () => {
    try {
      const history = await fetchChatHistory(access_token!);
      
      // Only replace messages if there's actual history
      if (history && history.length > 0) {
        setMessages(history);
      }
      // Otherwise, keep the initial welcome messages
    } catch (error) {
      console.error("Failed to load chat history:", error);
      // On error, we keep the initial messages too
    }
  };

  const handleSendMessage = async () => {
    if (!currentMessage.trim() || !access_token) return;

    setIsSending(true);
    
    try {
      const newMessage: ChatMessage = {
        id: Date.now().toString(),
        content: currentMessage,
        type: 'user',
        timestamp: new Date().toISOString()
      };

      setMessages(prev => [...prev, newMessage]);
      setCurrentMessage("");

      const response = await sendMessage(currentMessage, access_token);
      setMessages(prev => [...prev, response]);
    } catch (error) {
      console.error("Failed to send message:", error);
    } finally {
      setIsSending(false);
    }
  };

  // Scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  useEffect(() => {
    if (!access_token) {
      window.location.href = '/auth/login';
    }
  }, [access_token]);

  if (isPageLoading) {
    return <div className="flex items-center justify-center h-[600px]">Loading...</div>;
  }

  if (!access_token) {
    return null;
  }

  return (
    <div className="flex flex-col h-full relative max-w-3xl mx-auto">
      {/* Messages Container */}
      <div className="flex-1 space-y-8 overflow-y-auto absolute top-0 left-0 right-0 bottom-[100px] pr-2">
        <div className="px-6 space-y-4 pb-8">
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
                <div className="px-2 text-md text-foreground">
                  {message.content}
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