'use server'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;


// Create a new project
export async function createProject(projectData: any) {
  console.log(projectData);
    try {
      const response = await fetch(`${API_BASE_URL}/create_project`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(projectData),
      });
      console.log(JSON.stringify(projectData));
      if (!response.ok) {
        throw new Error('Failed to create project');
      }
      return await response.json();
    } catch (error) {
      console.error('Error creating project:', error);
      throw error;
    }
  }
