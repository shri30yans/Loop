// actions.ts
import { createClient } from '@/utils/supabase/client'; 
import {ProjectSectionType, ProjectType} from "../types"


export async function addProject(postData: ProjectType) {
console.log(postData)
const supabase= createClient();
  const { error } = await supabase
    .from('projects')
    .insert(postData);

  if (error) {
    console.error('Error adding post: ', error);
    return null;
  }

  return postData;
}