// actions.ts
import { createClient } from '@/utils/supabase/client'; 
import {ProjectType} from "../types"

export async function fetchProjectInfo(id:String){
    const supabase= createClient();
    const { data, error } = await supabase
        .from('projects')
        .select()
        .eq("id",id)
        .single()

    if (error) {
        console.error('Error fetching post info', error);
        return null;
    }
  return data as ProjectType;
}