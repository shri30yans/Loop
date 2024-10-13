"use client";
import { Select, SelectItem } from "@nextui-org/select";
import { useEffect, useState } from "react";
import { fetchProjects } from "./actions";
import ProjectCard from "@/components/projectcard";
import { Skeleton } from "@nextui-org/skeleton";
import { ProjectType } from "../types";
import { heading } from "@/components/primitives";

export default function FeedPage() {
  const types = [
    { key: "ongoing", label: "Ongoing" },
    { key: "completed", label: "Completed" },
  ];
  const sortByOptions = [
    { key: "CreatedAt", label: "Newest" },
    { key: "Title", label: "Title" },
  ];

  const [projects, setProjects] = useState<ProjectType[]>([]);
  const [isLoaded, setIsLoaded] = useState(false);
  const [selectedType, setSelectedType] = useState(types[0].key);
  const [selectedSortBy, setSelectedSortBy] = useState(sortByOptions[0].key);

  // Fetch projects when the component mounts or filters change
  useEffect(() => {
    fetchProjects(selectedType, selectedSortBy, "").then((fetchedProjects: any | null) => {
      if (fetchedProjects) {
        setProjects(fetchedProjects);
      }
    });
  }, [selectedType, selectedSortBy]);

  useEffect(() => {
    if (projects.length > 0) {
      setIsLoaded(true);
    }
  }, [projects]);

  return (
    <div className="space-y-4">
      <div className={heading()}>Projects Feed</div>
      <div className="flex gap-4">
        {/* Type Select */}
        <Select
          label="Type"
          selectedKeys={[selectedType]}
          onSelectionChange={(key) => setSelectedType(key as string)}
          className="w-40"
        >
          {types.map((type) => (
            <SelectItem key={type.key}>{type.label}</SelectItem>
          ))}
        </Select>

        {/* Sort By Select */}
        <Select
          label="Sort by"
          selectedKeys={[selectedSortBy]}
          onSelectionChange={(key) => setSelectedSortBy(key as string)}
          className="w-40"
        >
          {sortByOptions.map((sortOption) => (
            <SelectItem key={sortOption.key}>{sortOption.label}</SelectItem>
          ))}
        </Select>
      </div>

      {/* Projects List */}
      <div className="flex flex-wrap gap-4 w-full">
        {projects.map((project) => (
          <div className="w-1/4" key={project.id}>
            <a href={`/projectpage?id=${project.id}`}>
              <ProjectCard isLoaded={isLoaded} title={project.title} body={project.description} tags={project.tags} />
            </a>
          </div>
        ))}
      </div>

      {/* Loading Skeleton (Optional) */}
      {/* {projects.length === 0 && (
        <div className="flex gap-4 w-full">
          {Array.from({ length: 4 }).map((_, index) => (
            <div className="w-1/4 py-4" key={index}>
              <div className="space-y-3 p-8">
                <Skeleton isLoaded={isLoaded} className="rounded-lg h-36 w-50" />
                <Skeleton isLoaded={isLoaded} className="w-3/5 h-3 rounded-lg" />
                <Skeleton isLoaded={isLoaded} className="w-4/5 h-2 rounded-lg" />
                <Skeleton isLoaded={isLoaded} className="w-3/5 h-2 rounded-lg" />
              </div>
            </div>
          ))}
        </div>
      )} */}
    </div>
  );
}
