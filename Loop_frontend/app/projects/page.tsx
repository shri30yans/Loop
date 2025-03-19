"use client";
import { Select, SelectItem } from "@nextui-org/select";
import { useEffect, useState } from "react";
import { getAllProjects } from "./actions";
import ProjectCard from "@/components/ui/projectcard";
import { Skeleton } from "@nextui-org/skeleton";
import { ProjectType } from "../types";
import { useAuthStore } from "@/lib/auth/authStore";
import { Input } from "@nextui-org/input";
import { Button } from "@nextui-org/button";

export default function FeedPage() {
  const [projects, setProjects] = useState<ProjectType[]>([]);
  const [totalProjects, setTotalProjects] = useState<number>(0);
  const [isLoaded, setIsLoaded] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");

  const access_token = useAuthStore((state) => state.access_token);

  useEffect(() => {
    if (access_token) {
      getAllProjects(access_token).then((fetchedData: any | null) => {
        if (fetchedData) {
          console.log(fetchedData);
          setProjects(fetchedData.projects || []);
          setTotalProjects(fetchedData.total || 0);
        }
      });
    }
  }, [access_token]);

  useEffect(() => {
    setIsLoaded(true);
  }, [projects]);

  const handleSearch = async () => {
    try {
      console.log(access_token)
      if (access_token) {
        const fetchedData = await getAllProjects(access_token, searchQuery);
        console.log("fetched data:",fetchedData)
        setProjects(fetchedData.projects || []);
        setTotalProjects(fetchedData.total || 0); 
      }
    } catch (error) {
      console.error("Search failed:", error);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleSearch();
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex gap-4">
        <div className="flex gap-2 w-1/4">
          <Input
            type="text"
            placeholder="Search..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onKeyDown={handleKeyPress}
            className="w-full"
          />
          <Button onClick={handleSearch} variant="flat" color="primary">
            Search
          </Button>
        </div>
      </div>
      <div>
        {totalProjects > 0 ? (
          <p className="text-gray-600 text-md">
            Search results: {totalProjects} project(s) found
          </p>
        ) : (
          <p className="text-gray-600">No projects found.</p>
        )}
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 w-full">
        {projects.length > 0 ? (
          projects.map((project) => (
            <div key={project.id}>
              <a href={`/projectpage?id=${project.id}`}>
                <ProjectCard
                  isLoaded={isLoaded}
                  title={project.title}
                  body={project.description}
                  tags={project.tags}
                />
              </a>
            </div>
          ))
        ) : (
          searchQuery && (
            <div className="col-span-full text-center">
              <p className="text-gray-500">Try refining your search query.</p>
            </div>
          )
        )}
      </div>
    </div>
  );
}
