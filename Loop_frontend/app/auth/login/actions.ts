'use server';

<<<<<<< HEAD
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
=======
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

    // if (response.status == 401) {
    //   throw new Error('Invalid credentials');
    // }

    if (!response.ok) {
      console.log('Login failed:', response.status, response.statusText);
      throw new Error('Login failed');
    }

    console.log('Login successful:', response);

    const data = await response.json();
    return data;


  } catch (error) {
    console.error('Login failed:', error);
    throw new Error('Login failed');
  }
}

>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
