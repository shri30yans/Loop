"use client";
import { Select, SelectItem } from "@nextui-org/select";
import { useEffect, useState, useCallback } from "react";
import { getAllProjects } from "./actions";
import { NetworkError, TimeoutError } from "@/utils/errors";
import ProjectCard from "@/components/ui/projectcard";
import { Skeleton } from "@nextui-org/skeleton";
import { ProjectType } from "../types";
import { useAuth } from "@/lib/auth/authStore";
import { Input } from "@nextui-org/input";
import { Button } from "@nextui-org/button";
import { Pagination } from "@nextui-org/pagination";
import debounce from "lodash/debounce";

const ITEMS_PER_PAGE = 12;

export default function FeedPage() {
  const [projects, setProjects] = useState<ProjectType[]>([]);
  const [totalProjects, setTotalProjects] = useState<number>(0);
  const [isLoaded, setIsLoaded] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);

  const token = useAuth((state) => state.token);

  const fetchProjects = useCallback(async (query: string = "") => {
    if (!token) return;

    try {
      setIsLoading(true);
      setError(null);
      const fetchedData = await getAllProjects(token, query);
      setProjects(Array.isArray(fetchedData.projects) ? fetchedData.projects : []);
      setTotalProjects(typeof fetchedData.total === 'number' ? fetchedData.total : 0);
      setCurrentPage(1); // Reset to first page on new search
    } catch (error: unknown) {
      if (error instanceof NetworkError || error instanceof TimeoutError) {
        setError(error.message);
      } else if (error instanceof Error) {
        setError(error.message);
      } else {
        setError(query ? "Search failed. Please try again." : "An unexpected error occurred");
      }
      setProjects([]);
      setTotalProjects(0);
    } finally {
      setIsLoading(false);
    }
  }, [token]);

  // Debounced search function
  const debouncedSearch = useCallback(
    debounce((query: string) => {
      fetchProjects(query);
    }, 500),
    [fetchProjects]
  );

  useEffect(() => {
    if (searchQuery) {
      debouncedSearch(searchQuery);
    } else {
      fetchProjects();
    }
  }, [searchQuery, debouncedSearch]);

  useEffect(() => {
    setIsLoaded(true);
  }, [projects]);

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(e.target.value);
  };

  // Calculate pagination
  const totalPages = Math.ceil(totalProjects / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const endIndex = startIndex + ITEMS_PER_PAGE;
  const currentProjects = Array.isArray(projects) ? projects.slice(startIndex, endIndex) : [];

  return (
    <div className="space-y-4">
      <div className="flex gap-4">
        <div className="flex gap-2 w-1/4">
          <Input
            type="text"
            placeholder="Search..."
            value={searchQuery}
            onChange={handleSearchChange}
            className="w-full"
          />
        </div>
      </div>
      {error ? (
        <div className="error-container p-4 mb-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-600">{error}</p>
          <button
            onClick={() => fetchProjects(searchQuery)}
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

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 w-full">
        {currentProjects.length > 0
          ? currentProjects.map((project) => (
              <div key={project.id}>
                <a href={`/projectpage?id=${project.id}`}>
                  <ProjectCard
                    isLoaded={isLoaded}
                    title={project.title}
                    body={project.description}
                    tags={Array.isArray(project.tags) ? project.tags : []}
                  />
                </a>
              </div>
            ))
          : searchQuery && (
              <div className="col-span-full text-center">
                <p className="text-gray-500">Try refining your search query.</p>
              </div>
            )}
      </div>
      
      {totalPages > 1 && (
        <div className="flex justify-center mt-8">
          <Pagination
            total={totalPages}
            page={currentPage}
            onChange={setCurrentPage}
            showControls
          />
        </div>
      )}
    </div>
  );
}
