'use server';

import { api } from "@/utils/api";
import { UserType } from "../types";

export async function getUserInfo(token: string, id: string): Promise<UserType | null> {
  try {
    return await api.users.getProfile(token, id);
  } catch (error) {
    console.error("Error fetching user info:", error);
    return null;
  }
}
