'use client'
import { title,sectionheading} from "@/components/primitives";
import { Select, SelectItem} from "@nextui-org/select";
import { createClient } from '@/utils/supabase/server';
import { Card, CardBody, CardFooter, CardHeader } from "@nextui-org/card";
import {Image} from "@nextui-org/image";
import {Divider} from "@nextui-org/divider";
import { useEffect,useState } from "react";
import { Button } from "@nextui-org/button";
import { fetchPosts } from "./actions";
import {postCard} from "@/components/postcard";
import {Skeleton} from "@nextui-org/skeleton";

export default function FeedPage() {
  const type = [
    {key: "posts", label: "Posts"},
    {key: "projects", label: "Projects"},
  ];
  const sortby = [
    {key: "best", label: "Best"},
    {key: "new", label: "New"},
    {key: "top", label: "Top"},
    {key: "controversial", label: "Controversial"},
  ];

  const [posts, setPosts] = useState([]);
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    fetchPosts("","","").then((fetchedPosts) => {
      setPosts(fetchedPosts);
    });
  }, []);

  useEffect(() => {
    if (posts.length > 0) {
      setIsLoaded(true);
    }
  }, [posts]);


  return (
    <div className="space-y-4">
          {/* <div>
            <h1 className={title()}>Choose what you want to see</h1>
          </div> */}
          <div className="space-x-4 md-5">
            <Select
              label="Feed"
              selectionMode="multiple"
              className="w-40"
            >
              {type.map((data) => (
                <SelectItem key={data.key}>
                  {data.label}
                </SelectItem>
              ))}
            </Select>
            <Select
              label="Sort by"
              className="w-40"
            >
              {sortby.map((data) => (
                <SelectItem key={data.key}>
                  {data.label}
                </SelectItem>
              ))}
            </Select>
          </div>
          <div>
              {posts.map((data) => (
               postCard(data.title,data.body)
              ))}

              <div className="space-y-8">
                  {Array.from({ length: 3 }).map((_, index) => (
                      <div className="w-full flex flex-col gap-2 space-y-2">
                        <Skeleton isLoaded={isLoaded} className="h-6 w-1/2 rounded-lg"/>
                        <Skeleton isLoaded={isLoaded} className="h-4 w-4/5 rounded-lg"/>
                        <Skeleton isLoaded={isLoaded} className="h-4 w-4/5 rounded-lg"/>
                      </div>
                  ))}
              </div>

              
          </div>
    </div>
  );
}
