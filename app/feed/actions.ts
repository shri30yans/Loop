// actions.ts
import { createClient } from '@/utils/supabase/client'; 

export async function fetchPosts(type:any, sortBy:any,timeRange:any){
    const supabase= createClient();
    const { data, error } = await supabase
        .from('posts')
        .select('title,body')

    if (error) {
        console.error('Error adding post: ', error);
        return null;
    }
  return data;
}