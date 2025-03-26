'use server';

import { api } from "@/utils/api";

export async function deleteAccount(token: string): Promise<void> {
  try {
    return await api.users.deleteAccount(token);
  } catch (error) {
    console.error('Account deletion failed:', error);
    throw error;
  }
}