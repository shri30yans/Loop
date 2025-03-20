<<<<<<< HEAD
'use server';

import { api } from "@/utils/api";

export async function updatePassword(token: string, currentPassword: string, newPassword: string): Promise<void> {
  try {
    return await api.auth.changePassword(token, currentPassword, newPassword);
  } catch (error) {
    console.error('Password update failed:', error);
    throw error;
=======
// actions.ts
'use server'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

export async function updatePassword(token: string, currentPassword: string, newPassword: string) {
  try {
    const response = await fetch(`${API_BASE_URL}/auth/edit_password`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({currentPassword, newPassword }),
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
