'use server'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

// Fetch a specific project
export async function getProjectInfo(id: string) {
  try {
    const response = await fetch(`${API_BASE_URL}/project/fetch_project/?project-id=${id}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch project ${id}`);
    }
    return await response.json();
  } catch (error) {
    console.error(`Error fetching project ${id}:`, error);
    throw error;
  }
}