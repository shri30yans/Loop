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

    // Handle common HTTP status codes
    switch (response.status) {
      case 401:
        throw new Error('Unauthorized');
      case 404:
        return { success: false, data: null } as T;
      case 500:
        throw new Error('Internal server error');
    }

    if (!response.ok) {
      if (!response.status) {
        throw new NetworkError();
      }
      throw new Error(`Server error: ${response.status}`);
    }

    const data = await response.json();
    return data as T;
  } catch (error: unknown) {
    if (error instanceof Error) {
      if (error.name === 'AbortError') {
        throw new TimeoutError();
      }
      if (error instanceof NetworkError || error instanceof TimeoutError) {
        throw error;
      }
      console.error('API request failed:', error);
    }
    throw new Error('An unexpected error occurred. Please try again.');
  }
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
