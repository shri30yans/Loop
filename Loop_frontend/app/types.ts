export type AuthResponse = {
  access_token: string;
  user_id: string;
  expires_at: string;
};

export type ProjectsResponse = {
  projects: ProjectType[];
  total: number;
};

export type RegisterData = {
  email: string;
  password: string;
  username: string; 
};

export interface UserType {
    id: string;
    username: string;
    email: string;
    bio?: string;
    location?: string;
    avatar_url?: string;
    created_at: string;
    updated_at?: string;
    projects?: ProjectType[];
}

interface Comment {
    content: string;
    author: string;
    date: string;
  }

export interface PostType {
    title: string;
    body: string;
};

export interface ProjectType {
    id?: string; 
    title: string;
    description: string;
    status: string;
    introduction: string;
    sections: ProjectSectionType[];
    owner_id: string;
    tags: string[];
    comments?: Comment[];
    owner?: UserType;
};

export type ProjectSectionType = {
    title: string;
    content: string;
};
