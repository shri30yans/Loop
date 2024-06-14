// actions.ts
import { createClient } from '@/utils/supabase/client'; 

export async function addPost(postData: any) {
const supabase= createClient();
  const { error } = await supabase
    .from('posts')
    .insert(postData);

  if (error) {
    console.error('Error adding post: ', error);
    return null;
  }

  return postData;
}