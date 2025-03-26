"use client";
import { Card, CardBody } from "@nextui-org/card";
import { Image } from "@nextui-org/image";
import { Avatar, AvatarIcon } from "@nextui-org/avatar";
import { subheading, heading } from "@/components/ui/primitives";
import { ProjectType } from "../types";
import { useEffect, useState } from "react";
import { getProjectInfo } from "./actions";
import { Chip } from "@nextui-org/chip";
import { useAuthStore } from "@/lib/auth/authStore";
import { useSearchParams } from 'next/navigation';

export default function ProjectPage() {
  const [project, setProject] = useState<ProjectType | null>(null);
  const access_token = useAuthStore((state) => state.access_token);
  const searchParams = useSearchParams();
  const id = searchParams.get("id");
  
  useEffect(() => {
    if (access_token && id) {
      getProjectInfo(access_token, id).then(
        (fetchedProject: ProjectType | null) => {
          if (fetchedProject) {
            // Ensure tags is never undefined
            fetchedProject.tags = fetchedProject.tags || [];
            setProject(fetchedProject);
            console.log("Fetched project:", fetchedProject);
          }
        }
      );
    }
  }, [access_token, id]);

  if (!project) {
    return <div>Loading...</div>;
  }

  return (
    <div className="max-w-full">
      <div className="space-y-4">
        <div className="flex gap-2 h-full">
          <div>
            <Image
              width={800}
              height={600}
              alt="NextUI hero Image"
              className="p-2"
              src="https://www.liquidplanner.com/wp-content/uploads/2019/04/HiRes-17.jpg"
            />
          </div>

          <div className="p-6 max-w-2xl">
            <div className={`${heading({ size: "lg" })}`}>{project.title}</div>
            <div className="h-full pl-2 pt-1 flex flex-col space-y-6">
              <div className={`${subheading({ size: "lg" })}`}>
                {project.description}
              </div>
              <div className="wrap-text break-words">
                {project.introduction}
              </div>
              <div>
              <a href={`profile?id=${project.owner_id}`} className="flex items-center space-x-3">
                  <Avatar
                    icon={<AvatarIcon />}
                    classNames={{
                      base: "bg-gradient-to-br from-[#00B4DB] to-[#0083B0]",
                      icon: "text-black/80",
                    }}
                    className="w-12 h-12"
                    alt={project.owner_id}
                  />
                  <span className="text-base font-medium">
                    {project.owner?.username}
                  </span>
                </a>
               </div>
              <div className="mt-auto space-y-4">
                <div className="flex flex-wrap gap-2">
                  {project.tags?.map((tag, index) => (
                    <Chip key={index} size="sm" radius="sm" variant="bordered">
                      {tag}
                    </Chip>
                  ))}
                </div>

              </div>
            </div>
          </div>
        </div>

        <div>
          {project.sections?.map((card, index) => (
            <div
              className="flex w-full flex-col space-y-6"
              key={index}
            >
              <div className="space-y-2">
                <div className="w-full space-y-2 px-6 pb-6 pt-2">
                  <div className={subheading({ size: "lg" })}>{card.title}</div>
                  <div className="wrap-text break-words">{card.content}</div>
                </div>
              </div>
            </div>
          ))}
        </div>
        <div></div>
      </div>
    </div>
  );
}
