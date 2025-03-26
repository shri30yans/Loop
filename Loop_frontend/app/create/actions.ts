'use server'
import { api } from "@/utils/api";
import { ProjectType } from "../types";

export async function createProject(access_token: string, project: ProjectType): Promise<ProjectType> {
  return api.projects.create(access_token, project);
}
