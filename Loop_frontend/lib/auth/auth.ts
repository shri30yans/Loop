'use server';

import { useAuthStore } from './authStore';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

export async function login(email: string, password: string) {
  try {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      throw new Error('Login failed');
    }

    const data = await response.json();

    // Parse the expires_at date
    const expiresAt = new Date(data.expires_at).toUTCString();

    // Set refresh token in an HTTP-only cookie with expiration
    document.cookie = `refresh_token=${data.refresh_token}; path=/; HttpOnly; Secure; expires=${expiresAt}`;

    // Set the auth state
    useAuthStore.getState().setAuth(data.refresh_token, data.user_id, data.expires_at);

    return data;
  } catch (error) {
    throw new Error('Login failed');
  }
}
