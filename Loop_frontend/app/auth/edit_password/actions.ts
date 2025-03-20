'use server';

import { api } from "@/utils/api";

export async function updatePassword(token: string, currentPassword: string, newPassword: string): Promise<void> {
  try {
    return await api.auth.changePassword(token, currentPassword, newPassword);
  } catch (error) {
    console.error('Password update failed:', error);
    throw error;
  }
}
