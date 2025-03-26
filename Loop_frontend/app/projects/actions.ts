'use server'

import { api } from '@/utils/api';
import { ProjectsResponse } from '../types';

export async function getAllProjects(access_token: string, searchKeyword?: string): Promise<ProjectsResponse> {
  return api.projects.searchProject(access_token, searchKeyword);
}
