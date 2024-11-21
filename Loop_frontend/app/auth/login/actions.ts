'use server';

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

