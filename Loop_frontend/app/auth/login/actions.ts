'use server';

import { api } from "@/utils/api";
import { AuthResponse } from "@/app/types";

export async function login(email: string, password: string): Promise<AuthResponse> {
  try {
    return await api.auth.login(email, password);
  } catch (error) {
    console.error('Login failed:', error);
    throw error;
  }
}
