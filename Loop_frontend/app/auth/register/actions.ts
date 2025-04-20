'use server';

import { api } from "@/utils/api";
import { AuthResponse } from "@/app/types";

export async function register(username: string, email: string, password: string): Promise<AuthResponse> {
  try {
    console.log('Register action called with:', { username, email, password });
    const registerData = {
      username: username,
      email: email,
      password: password
    };
    return await api.auth.register(registerData);
  } catch (error) {
    console.error('Registration failed:', error);
    throw error;
  }
}
