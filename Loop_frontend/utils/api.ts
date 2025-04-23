import { NetworkError, TimeoutError } from './errors';
import { ProjectType, UserType, AuthResponse, ProjectsResponse, RegisterData } from '@/app/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;
const DEFAULT_TIMEOUT = 5000;

type RequestConfig = {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  body?: any;
  timeout?: number;
  headers?: Record<string, string>;
};

function toCurlCommand(url: string, method: string, headers: Record<string, string>, body?: any): string {
  let curl = `curl -X ${method} '${url}'`;
  
  // Add headers
  Object.entries(headers).forEach(([key, value]) => {
    curl += ` \\\n  -H '${key}: ${value}'`;
  });

  // Add body if present
  if (body) {
    curl += ` \\\n  -d '${JSON.stringify(body)}'`;
  }

  return curl;
}

export async function apiRequest<T>(
  endpoint: string,
  access_token: string,
  config: RequestConfig = {}
): Promise<T> {
  const {
    method = 'GET',
    body,
    timeout = DEFAULT_TIMEOUT,
    headers = {},
  } = config;

  const url = `${API_BASE_URL}${endpoint}`;
  const requestHeaders = {
    'Authorization': `Bearer ${access_token}`,
    'Content-Type': 'application/json',
    ...headers,
  };

  // Log equivalent curl command for debugging
  console.log('\nEquivalent curl command:');
  console.log(toCurlCommand(url, method, requestHeaders, body));
  console.log(); // Empty line for readability

  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeout);

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      method,
      headers: {
        'Authorization': `Bearer ${access_token}`,
        'Content-Type': 'application/json',
        ...headers,
      },
      body: body ? JSON.stringify(body) : undefined,
      credentials: 'include',
      mode: 'cors',
      signal: controller.signal,
    });

    clearTimeout(timeoutId);

    // Check if we actually got a response from the server
    if (!response.status) {
      throw new NetworkError('Unable to connect to the server. Please check your connection and try again.');
    }

    // For any error response, try to get the response body
    if (!response.ok) {
      let errorData;
      let errorBody = '';
      try {
        errorBody = await response.text();
        console.log('Error response body:', errorBody); // Log raw error response
        try {
          errorData = JSON.parse(errorBody);
          console.log('Parsed error data:', errorData); // Log parsed error data
        } catch {
          // If JSON parsing fails, use the raw text
          errorData = { message: errorBody };
        }
      } catch {
        errorData = { message: `Server error: ${response.status}` };
      }

      // Enhanced error message with status code and response body
      let errorMessage = `Server Error (${response.status})`;
      if (errorData.message) {
        errorMessage += `: ${errorData.message}`;
      }
      if (errorBody && errorMessage !== errorBody) {
        errorMessage += `\nResponse: ${errorBody}`;
      }

      // If there are validation details, include them
      if (errorData.details) {
        const details = Object.entries(errorData.details)
          .map(([field, error]) => `${field}: ${error}`)
          .join(', ');
        errorMessage += `\nDetails: ${details}`;
      }

      throw new NetworkError(errorMessage, response, errorData);
    }

    let data = await response.json();
    console.log('Raw server response:', data); // Debug log
    
    // Handle development environment response structure
    if (Array.isArray(data) && data.length === 2) {
      if (data[0]?.[0] === '$@1' && data[1]) {
        // Extract the actual response data from the second element
        data = data[1];
      }
    }
    
    return data as T;
  } catch (error: unknown) {
    if (error instanceof Error) {
      // Handle timeout errors
      if (error.name === 'AbortError') {
        throw new TimeoutError();
      }
      
      // Handle network errors
      if (error.message.includes('Failed to fetch') || error.message.includes('Network request failed')) {
        throw new NetworkError('Unable to connect to the server. Please check your connection and try again.');
      }
      
      // Check online status only in browser environment
      if (typeof window !== 'undefined' && !window.navigator.onLine) {
        throw new NetworkError('You appear to be offline. Please check your internet connection.');
      }

      // If it's already a NetworkError or TimeoutError, just rethrow it
      if (error instanceof NetworkError || error instanceof TimeoutError) {
        throw error;
      }

      console.error('API request failed:', error);
      throw error;
    }
    throw new Error('An unexpected error occurred. Please try again.');
  }
}

// Add ChatMessage interface
export interface ChatMessage {
  id: string;
  content: string;
  type: 'user' | 'llm';
  timestamp: string;
}

// Type-safe API functions
export const api = {
  auth: {
    login: async (email: string, password: string): Promise<AuthResponse> => {
      return apiRequest('/auth/login', '', {
        method: 'POST',
        body: { email, password },
      });
    },
    register: async (data: RegisterData): Promise<AuthResponse> => {
      return apiRequest('/auth/register', '', {
        method: 'POST',
        body: data,
      });
    },
    changePassword: async (access_token: string, currentPassword: string, newPassword: string): Promise<void> => {
      return apiRequest('/auth/edit_password', access_token, {
        method: 'PUT',
        body: { currentPassword, newPassword },
      });
    },
  },
  projects: {
    searchProject: async (access_token: string, searchKeyword?: string): Promise<ProjectsResponse> => {
      const queryParams = searchKeyword ? `?keyword=${encodeURIComponent(searchKeyword)}` : '';
      const response = await apiRequest(`/project/search${queryParams}`, access_token);
      return response as ProjectsResponse;
    },
    getProject2w: async (access_token: string, searchKeyword?: string): Promise<ProjectsResponse> => {
      const queryParams = searchKeyword ? `?keyword=${encodeURIComponent(searchKeyword)}` : '';
      const response = await apiRequest(`/project/get${queryParams}`, access_token);
      return response as ProjectsResponse;
    },
    getProject: async (access_token: string, id: string): Promise<ProjectType> => {
      const response = await apiRequest(`/project/${id}`, access_token);
      console.log(response);
      return response as ProjectType;
    },
    create: async (access_token: string, data: Partial<ProjectType>): Promise<ProjectType> => {
      return apiRequest('/project', access_token, {
        method: 'POST',
        body: data,
      });
    },
    update: async (access_token: string, id: string, data: Partial<ProjectType>): Promise<ProjectType> => {
      return apiRequest(`/project/${id}`, access_token, {
        method: 'PUT',
        body: data,
      });
    },
  },
  users: {
    getProfile: async (access_token: string, userId: string): Promise<UserType> => {
      return apiRequest(`/user/${userId}`, access_token);
    },
    update: async (access_token: string, data: Partial<UserType>): Promise<UserType> => {
      return apiRequest('/user/update', access_token, {
        method: 'PUT',
        body: data,
      });
    },
    
    deleteAccount: async (access_token: string): Promise<void> => {
      return apiRequest('/user/delete', access_token, {
        method: 'PUT',
      });
    },
  },
};

// Chat API Functions
export const chatApi = {
  sendMessage: async (message: string, token: string): Promise<ChatMessage> => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/chat/conversation`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ message })
    });

    if (!response.ok) {
      throw new Error('Failed to send message');
    }

    return response.json();
  },

  fetchChatHistory: async (token: string): Promise<ChatMessage[]> => {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/chat/history`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    if (!response.ok) {
      throw new Error('Failed to fetch chat history');
    }

    return response.json();
  }
};
