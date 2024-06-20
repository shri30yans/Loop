"use client";
import { Select, SelectItem } from "@nextui-org/select";
import { useEffect, useState } from "react";
import { fetchProjects } from "./actions";
import ProjectCard from "@/components/projectcard";
import { Skeleton } from "@nextui-org/skeleton";
import  {ProjectType}  from "../types";

export default function FeedPage() {
  const type = [
    { key: "posts", label: "Posts" },
    { key: "projects", label: "Projects" },
  ];
  const sortby = [
    { key: "best", label: "Best" },
    { key: "new", label: "New" },
    { key: "top", label: "Top" },
    { key: "controversial", label: "Controversial" },
  ];

  const [posts, setPosts] = useState<ProjectType[]>([]);
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    fetchProjects("", "", "").then((fetchedProjects:ProjectType[] | null) => {
      if (fetchedProjects) {
        setPosts(fetchedProjects);
      }
    });
  }, []);

  useEffect(() => {
    if (posts.length > 0) {
      setIsLoaded(true);
    }
  }, [posts]);

  return (
    <div className="space-y-4">
      <div className="space-x-4 md-5">
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
        {posts.map((data) => (
          <div className="w-1/4">
            <ProjectCard isLoaded= {isLoaded} title={data.title} body={data.body} category={"AI"} />
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
