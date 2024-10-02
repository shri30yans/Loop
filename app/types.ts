import { ProjectReference } from "typescript";

export interface PostType {
    title: string;
    body: string;
};

export interface ProjectType {
    id?:string;
    title: string;
    description: string;
    introduction: string;
    sections: ProjectSectionType[];
    user: string;
    tags?: string;
};

export type ProjectSectionType = {
    id: string;
    title: string;
    body: string;
};