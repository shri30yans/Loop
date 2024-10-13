import { useEffect, useState } from "react";
import { fetchPosts } from "./actions";
import PostCard from "@/components/postcard";
import { Select, SelectItem } from "@nextui-org/select";
import { PostType } from "../types";
import { Skeleton } from "@nextui-org/skeleton";

export default function FeedPage() {
  const types = [
    { key: "posts", label: "Posts" },
    { key: "projects", label: "Projects" },
  ];

  const sortby = [
    { key: "best", label: "Best" },
    { key: "new", label: "New" },
    { key: "top", label: "Top" },
    { key: "controversial", label: "Controversial" },
  ];

  const [posts, setPosts] = useState<PostType[]>([]);
  const [isLoaded, setIsLoaded] = useState(false);
  const [selectedType, setSelectedType] = useState("posts");
  const [selectedSortBy, setSelectedSortBy] = useState("best");

  useEffect(() => {
    fetchPosts(selectedType, selectedSortBy, "").then((fetchedPosts: PostType[] | null) => {
      if (fetchedPosts) {
        setPosts(fetchedPosts);
      }
    });
  }, [selectedType, selectedSortBy]);

  useEffect(() => {
    if (posts.length > 0) {
      setIsLoaded(true);
    }
  }, [posts]);

  return (
    <div className="space-y-4">
      <div className="space-x-4 md-5">
        <Select
          label="Sort by"
          className="w-40"
          onChange={(value) => setSelectedSortBy(value)}
        >
          {sortby.map((data) => (
            <SelectItem key={data.key} value={data.key}>
              {data.label}
            </SelectItem>
          ))}
        </Select>

        <Select
          label="Type"
          className="w-40"
          onChange={(value) => setSelectedType(value)}
        >
          {types.map((data) => (
            <SelectItem key={data.key} value={data.key}>
              {data.label}
            </SelectItem>
          ))}
        </Select>
      </div>

      <div>
        {posts.map((data) => (
          <PostCard key={data.id} title={data.title} body={data.body} />
        ))}

        {/* Skeleton for loading */}
        <div className="space-y-8">
          {Array.from({ length: 3 }).map((_, index) => (
            <div key={index} className="w-full flex flex-col gap-2 space-y-2">
              <Skeleton isLoaded={isLoaded} className="h-6 w-1/2 rounded-lg" />
              <Skeleton isLoaded={isLoaded} className="h-4 w-4/5 rounded-lg" />
              <Skeleton isLoaded={isLoaded} className="h-4 w-4/5 rounded-lg" />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
