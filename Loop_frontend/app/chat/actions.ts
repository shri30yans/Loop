import { ChatMessage, chatApi } from '../../utils/api';

// Re-export the ChatMessage type for components that import from here
export type { ChatMessage };

export async function sendMessage(message: string, access_token: string): Promise<ChatMessage> {
  return chatApi.sendMessage(message, access_token);
}

export async function fetchChatHistory(access_token: string): Promise<ChatMessage[]> {
  return chatApi.fetchChatHistory(access_token);
}
