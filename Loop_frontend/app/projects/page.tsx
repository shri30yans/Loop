"use client";
import { Select, SelectItem } from "@nextui-org/select";
import { useEffect, useState } from "react";
import { getAllProjects } from "./actions";
import ProjectCard from "@/components/ui/projectcard";
import { Skeleton } from "@nextui-org/skeleton";
import  {ProjectType}  from "../types";

export default function FeedPage() {
  const type = [
    { key: "projects", label: "Posts" },
    { key: "projects", label: "Projects" },
  ];
  const sortby = [
    { key: "best", label: "Best" },
    { key: "new", label: "New" },
    { key: "top", label: "Top" },
    { key: "controversial", label: "Controversial" },
  ];

  const [projects, setProjects] = useState<ProjectType[]>([]);
  const [isLoaded, setIsLoaded] = useState(false);


  useEffect(() => {
    getAllProjects().then((fetchedProjects: any | null) => {
      if (fetchedProjects) {
        setProjects(fetchedProjects);
      }
    });
  }, []);

  useEffect(() => {
    if (projects.length > 0) {
      setIsLoaded(true);
    }
  }, [projects]);

  return (
    <div className="space-y-4">
      <div className="flex gap-4">
        <Select label="Feed" selectionMode="multiple" className="w-40">
          {type.map((data) => (
            <SelectItem key={data.key}>{data.label}</SelectItem>
          ))}
        </Select>
        <Select label="Sort by" className="w-40">
          {sortby.map((data) => (
            <SelectItem key={data.key}>{data.label}</SelectItem>
          ))}
        </Select>
      </div>

      <div className="flex flex-wrap gap-4 w-full">
        {projects.map((project) => (
          <div className="w-1/4">
            <a href={`/projectpage?id=${project.}`}>
            <ProjectCard isLoaded={isLoaded} title={project.title} body={project.description} tags={project.tags} />
            </a>
          </div>
        ))}
      </div>
      {/* <div className="flex gap-4 w-full">
        {Array.from({ length: 4 }).map((_, index) => (
          <div className="w-1/4 py-4">
            <div className="space-y-3 p-8">
              <Skeleton isLoaded={isLoaded} className="rounded-lg h-36 w-50" />
              <Skeleton isLoaded={isLoaded} className="w-3/5 h-3 rounded-lg" />
              <Skeleton isLoaded={isLoaded} className="w-4/5 h-2 rounded-lg" />
              <Skeleton isLoaded={isLoaded} className="w-3/5 h-2 rounded-lg" />
            </div>
          </div>
        ))}
      </div> */}
    </div>
  );
}
