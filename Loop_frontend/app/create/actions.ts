'use server'
import { ProjectSectionType, ProjectType } from "../types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;


// Create a new project
export async function createProject(project : ProjectType) {
  console.log(project);
    try {
      const response = await fetch(`${API_BASE_URL}/project/create_project`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({...project}),
      });
      console.log(JSON.stringify(project));
      if (!response.ok) {
        console.log(response);
        throw new Error('Failed to create project');
      }
      return await response.json();
    } catch (error) {
      console.error('Error creating project:', error);
      throw error;
    }
  }
