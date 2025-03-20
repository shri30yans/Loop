<<<<<<< HEAD
'use server';

import { api } from "@/utils/api";

export async function deleteAccount(token: string): Promise<void> {
  try {
    return await api.users.deleteAccount(token);
  } catch (error) {
    console.error('Account deletion failed:', error);
    throw error;
=======
'use server'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

export async function deleteAccount(token: string) {
  try {
    const response = await fetch(`${API_BASE_URL}/user/delete_account`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    });

    if (!response.ok) {
      if (response.status === 401) {
        throw new Error('Unauthorized');
      }
      throw new Error('Failed to fetch user info');
    }

    if (!response.ok) {
      throw new Error('Failed to update password');
    }
    return await response.json();

  } catch (error) {
    throw new Error('Password update failed');
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
  }
}
