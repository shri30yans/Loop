"use client";
import { Select, SelectItem } from "@nextui-org/select";
import { useEffect, useState } from "react";
import { getAllProjects } from "./actions";
<<<<<<< HEAD
import { NetworkError, TimeoutError } from "@/utils/errors";
=======
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
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
<<<<<<< HEAD
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const access_token = useAuthStore((state) => state.access_token);

  const fetchProjects = async () => {
    if (!access_token) return;
    
    try {
      setIsLoading(true);
      setError(null);
      const fetchedData = await getAllProjects(access_token);
      setProjects(fetchedData.projects || []);
      setTotalProjects(fetchedData.total || 0);
    } catch (error: unknown) {
      if (error instanceof NetworkError || error instanceof TimeoutError) {
        setError(error.message);
      } else if (error instanceof Error) {
        setError(error.message);
      } else {
        setError('An unexpected error occurred');
      }
      setProjects([]);
      setTotalProjects(0);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchProjects();
=======

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
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
  }, [access_token]);

  useEffect(() => {
    setIsLoaded(true);
  }, [projects]);

  const handleSearch = async () => {
<<<<<<< HEAD
    if (!access_token) return;
    
    try {
      setIsLoading(true);
      setError(null);
      const fetchedData = await getAllProjects(access_token, searchQuery);
      setProjects(fetchedData.projects || []);
      setTotalProjects(fetchedData.total || 0);
    } catch (error: unknown) {
      if (error instanceof NetworkError || error instanceof TimeoutError) {
        setError(error.message);
      } else if (error instanceof Error) {
        setError(error.message);
      } else {
        setError('Search failed. Please try again.');
      }
      setProjects([]);
      setTotalProjects(0);
    } finally {
      setIsLoading(false);
=======
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
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
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
<<<<<<< HEAD
      {error ? (
        <div className="error-container p-4 mb-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-600">{error}</p>
          <button 
            onClick={handleSearch} 
            className="mt-2 px-4 py-2 bg-red-100 text-red-700 rounded hover:bg-red-200"
          >
            Retry
          </button>
        </div>
      ) : (
        <div>
          {isLoading ? (
            <p className="text-gray-600">Loading projects...</p>
          ) : totalProjects > 0 ? (
            <p className="text-gray-600 text-md">
              Search results: {totalProjects} project(s) found
            </p>
          ) : (
            <p className="text-gray-600">No projects found.</p>
          )}
        </div>
      )}
=======
      <div>
        {totalProjects > 0 ? (
          <p className="text-gray-600 text-md">
            Search results: {totalProjects} project(s) found
          </p>
        ) : (
          <p className="text-gray-600">No projects found.</p>
        )}
      </div>
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7

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
