'use server';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;


export async function register(name:string, email: string, password: string) {
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

    if (!response.ok) {
      throw new Error('Registration failed');
    }

    const data = await response.json();
    return data;
  } catch (error) {
    throw new Error('Registration failed');
  }
}

