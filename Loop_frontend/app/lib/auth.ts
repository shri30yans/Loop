// lib/auth.ts
import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';

interface User {
  id: number;
  email: string;
  created_at: string;
}

interface AuthState {
  token: string | null;
  user: User | null;
  setAuth: (token: string | null, user: User | null) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      setAuth: (token, user) => set({ token, user }),
      logout: () => set({ token: null, user: null }),
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => sessionStorage), // or localStorage
    }
  )
);

// Custom fetch wrapper with authentication
export async function fetchWithAuth(
  url: string,
  options: RequestInit = {}
): Promise<Response> {
  const token = useAuthStore.getState().token;
  
  const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const fullUrl = `${baseUrl}${url}`;

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(fullUrl, {
    ...options,
    headers,
  });

  if (response.status === 401) {
    // Handle unauthorized error - clear auth and redirect
    useAuthStore.getState().logout();
    window.location.href = '/login';
  }

  return response;
}

// Server-side auth utilities
export async function validateToken(token: string | null): Promise<boolean> {
  if (!token) return false;

  try {
    const response = await fetch(`${process.env.API_URL}/auth/validate`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.ok;
  } catch (error) {
    return false;
  }
}

// Authentication API functions
export const authApi = {
  login: async (email: string, password: string) => {
    const response = await fetchWithAuth('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      throw new Error('Login failed');
    }

    return response.json();
  },

  register: async (email: string, password: string) => {
    const response = await fetchWithAuth('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      throw new Error('Registration failed');
    }

    return response.json();
  },

  logout: async () => {
    useAuthStore.getState().logout();
  },

  // Example of a protected API call
  getProfile: async () => {
    const response = await fetchWithAuth('/profile', {
      method: 'GET',
    });

    if (!response.ok) {
      throw new Error('Failed to fetch profile');
    }

    return response.json();
  },
};

// Type guard for checking authentication
export function isAuthenticated(): boolean {
  const { token, user } = useAuthStore.getState();
  return !!token && !!user;
}

// Custom hook for protected routes
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export function useRequireAuth() {
  const router = useRouter();
  const { token, user } = useAuthStore();

  useEffect(() => {
    if (!token || !user) {
      router.push('/login');
    }
  }, [token, user, router]);

  return { user, isLoading: !token || !user };
}

// Error handling types
export interface ApiError {
  message: string;
  code?: string;
  status?: number;
}

// API response types
export interface AuthResponse {
  token: string;
  user: User;
}

// Auth form data types
export interface LoginFormData {
  email: string;
  password: string;
}

export interface RegisterFormData extends LoginFormData {
  confirmPassword: string;
}