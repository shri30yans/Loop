'use server';

import { ProjectType, UserType } from "../types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

export async function getUserInfo(token: string, id: string): Promise<UserType | null> {
  try {
    const response = await fetch(`${API_BASE_URL}/users/info?user_id=${id}`, {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      mode: "cors",
    });

    if (response.status === 401) {
      throw new Error("Unauthorized");
    }
    if (!response.ok) {
      throw new Error("Failed to fetch user info");
    }

    const data = await response.json();
    return {
      id: data.id,
      name: data.name,
      email: data.email,
      bio: data.bio || "No bio available",
      location: data.location || "Unknown",
      avatar_url: data.avatar_url || "https://i.pravatar.cc/150",
      created_at: data.created_at,
      updated_at: data.updated_at,
      projects: data.projects?.map((project: any) => ({
        project_id: project.project_id,
        title: project.title || "Untitled Project",
        description: project.description || "No description available",
        introduction: project.introduction,
        sections: project.sections || [],
        owner_id: project.owner_id,
        tags: project.tags || [],
        image_url: "https://via.placeholder.com/150",
        created_at: project.created_at,
        updated_at: project.updated_at,
      })),
    };
  } catch (error) {
    console.error("Error fetching user info:", error);
    return null;
  }
}
