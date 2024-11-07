export interface PostType {
    title: string;
    body: string;
};

export interface ProjectType {
    id: string;
    title: string;
    description: string;
    introduction: string;
    sections: ProjectSectionType[];
    //owner_id: string;
    tags: string;
};

export type ProjectSectionType = {
    update_number: number;
    title: string;
    body: string;
};