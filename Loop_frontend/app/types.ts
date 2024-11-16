export interface User {
    user_id: string;
    name: string;
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
};

export type ProjectSectionType = {
    section_number: number;
    title: string;
    body: string;
};