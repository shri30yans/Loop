'use server';

import { api } from "@/utils/api";
import { AuthResponse } from "@/app/types";

export async function register(username: string, email: string, password: string): Promise<AuthResponse> {
  try {
    return await api.auth.register({ username, email, password });
  } catch (error) {
    console.error('Registration failed:', error);
    throw error;
  }
}
