'use server';

<<<<<<< HEAD
import { api } from "@/utils/api";
import { AuthResponse } from "@/app/types";

export async function register(username: string, email: string, password: string): Promise<AuthResponse> {
  try {
    return await api.auth.register({ username, email, password });
  } catch (error) {
    console.error('Registration failed:', error);
    throw error;
=======
const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;


export async function register(name:string, email: string, password: string) {
  console.log(API_BASE_URL);
  try {
    const response = await fetch(`${API_BASE_URL}/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({name, email, password }),
    });
    console.log(name, email, password);
    console.log(response);

    if (response.status == 409) {
      throw new Error('User already exists');
    }


    else if (!response.ok) {
    console.log(response);

    throw new Error('Registration failed');
    }
    const data = await response.json();
    return data;
    
  } catch (error) {
    console.error('Error details:', error);
    throw new Error('Registration failed');
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
  }
}
