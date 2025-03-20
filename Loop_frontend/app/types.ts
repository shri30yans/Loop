<<<<<<< HEAD
export type AuthResponse = {
  access_token: string;
  user: UserType;
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
=======
export interface UserType {
    id: string;
    name: string;
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
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
<<<<<<< HEAD
    id?: string; 
=======
    id?: string;  // Changed from project_id to match backend JSON tag
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
    title: string;
    description: string;
    status: string;
    introduction: string;
    sections: ProjectSectionType[];
    owner_id: string;
    tags: string[];
    comments?: Comment[];
    owner?: UserType;
<<<<<<< HEAD
    created_at: string;
    updated_at: string;
};

export type ProjectSectionType = {
    title: string;
    body: string;
};
=======
};

export type ProjectSectionType = {
    index : number;
    title: string;
    content: string;
};
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
