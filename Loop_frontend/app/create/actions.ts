'use server'
import { ProjectSectionType, ProjectType } from "../types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

export async function createProject(refresh_token: string, project: ProjectType) {
  const projectData = {
    title: project.title,
    description: project.description,
    introduction: project.introduction,
    owner_id: parseInt(project.owner_id),
    tags: Array.isArray(project.tags) ? project.tags : [],
    sections: project.sections.map(section => ({
      title: section.title, 
      body: section.body,   
      section_number: section.section_number
    }))
  };

  console.log('Sending to backend:', JSON.stringify(projectData));
  
  try {
    const response = await fetch(`${API_BASE_URL}/project/create_project`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${refresh_token}`,
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      mode: 'cors',
      body: JSON.stringify(projectData)
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
