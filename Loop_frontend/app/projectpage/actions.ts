'use server'

import { ProjectType } from "../types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

function mapToProjectType(data: any): ProjectType {
  const projectData = Array.isArray(data) ? data[0] : data;

  // // Split tags string if it's a single string
  // const formattedTags = Array.isArray(projectData.tags) 
  //   ? projectData.tags[0].split(', ')
  //   : [];

  return {
    project_id: projectData.project_id.toString(),
    title: projectData.title,
    description: projectData.description,
    introduction: projectData.introduction,
    owner_id: projectData.owner_id.toString(),
    tags: projectData.tags,
    sections: projectData.sections.map((section: any) => ({
      section_number: section.section_number,
      title: section.title,
      body: section.body
    }))
  };
}

export async function getProjectInfo(refresh_token: string, id: string) {
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
    return mapToProjectType(data);
  } catch (error) {
    console.error('Error fetching projects:', error);
    throw error;
  }
}