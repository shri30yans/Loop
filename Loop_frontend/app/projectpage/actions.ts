'use server'

import { api } from "@/utils/api";
import { ProjectType } from "../types";

export async function getProjectInfo(access_token: string, id: string): Promise<ProjectType> {
  return api.projects.getProject(access_token, id);
}
