'use server'
import { ProjectSectionType, ProjectType } from "../types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

export async function createProject(refresh_token: string, project: ProjectType) {  
  try {
    const response = await fetch(`${API_BASE_URL}/project/create_project`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${refresh_token}`,
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      mode: 'cors',
      body: JSON.stringify(project)
    });

    if (response.status === 400) {
      const errorData = await response.text();
      console.error('Bad request details:', errorData);
      throw new Error('Invalid project data');
    }
    
    if (!response.ok) {
      const errorData = await response.text();
      console.log(response)
      throw new Error('Failed to create project');
    }

    return await response.json();

  } catch (error) {
    console.error('Error creating project:', error);
    throw error;
  }
}
