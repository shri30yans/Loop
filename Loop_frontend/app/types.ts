export interface UserType {
    id: string;
    name: string;
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
    project_id?: string;
    title: string;
    description: string;
    introduction: string;
    sections: ProjectSectionType[];
    owner_id: string;
    tags: string[];
    comments?: Comment[];
    owner?: UserType;
};


export type ProjectSectionType = {
    index: number;
    title: string;
    content: string;
};