// actions.ts
import { createClient } from '@/utils/supabase/client'; 


export async function fetchProjects(type:string, sortBy:string,timeRange:string){
    const supabase= createClient();
    const { data, error } = await supabase
        .from('projects')
        .select("id,title,description,tags")

    if (error) {
        console.error('Error adding post: ', error);
        return null;
    }
  return data;
}