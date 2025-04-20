export interface ChatMessage {
  id: string;
  content: string;
  type: 'user' | 'llm';
  timestamp: string;
}

export async function sendMessage(message: string, access_token: string): Promise<ChatMessage> {
  const response = await fetch('/api/chat/send', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${access_token}`
    },
    body: JSON.stringify({ message })
  });

  if (!response.ok) {
    throw new Error('Failed to send message');
  }

  return response.json();
}

export async function fetchChatHistory(access_token: string): Promise<ChatMessage[]> {
  const response = await fetch('/api/chat/history', {
    headers: {
      'Authorization': `Bearer ${access_token}`
    }
  });

  if (!response.ok) {
    throw new Error('Failed to fetch chat history');
  }

  return response.json();
}
