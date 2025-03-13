'use server';

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
  }
}
