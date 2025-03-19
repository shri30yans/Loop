'use server'
import { ProjectSectionType, ProjectType } from "../types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

export async function createProject(access_token: string, project: ProjectType) {  
  try {
    const response = await fetch(`${API_BASE_URL}/project/create`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${access_token}`,
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
      console.log('Error response:', response);
      console.log('Error data:', errorData);
      throw new Error('Failed to create project');
    }

    const data = await response.json();
    console.log('Project creation response:', data);
    return data;

  } catch (error) {
    console.error('Error creating project:', error);
    throw error;
  }
}
