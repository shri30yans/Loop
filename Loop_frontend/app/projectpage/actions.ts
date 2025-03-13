'use server'

import { ProjectType } from "../types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

export async function getProjectInfo(refresh_token: string, id: string): Promise<ProjectType> {
  try {
    const response = await fetch(`${API_BASE_URL}/project/get_project_info?project-id=${id}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${refresh_token}`,
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      mode: 'cors'
    });

    if (response.status === 401) {
      throw new Error('Unauthorized');
    }
    else if (!response.ok) {
      console.log(response)
      throw new Error('Failed to fetch projects');
    }
    
    const data = await response.json();
    // Return first item if response is an array
    return Array.isArray(data) ? data[0] : data;
  } catch (error) {
    console.error('Error fetching projects:', error);
    throw error;
  }
}