'use server'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;


// Create a new project
export async function createProject(projectData: any) {
    try {
      const response = await fetch(`${API_BASE_URL}/api/projects`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(projectData),
      });
      if (!response.ok) {
        throw new Error('Failed to create project');
      }
      return await response.json();
    } catch (error) {
      console.error('Error creating project:', error);
      throw error;
    }
  }