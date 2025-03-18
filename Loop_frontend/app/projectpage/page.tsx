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
import { useSearchParams } from "next/navigation";
import { getRandomBackground } from "@/utils/randomimage";

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
            setProject(fetchedProject);
            console.log(fetchedProject);
            console.log(fetchedProject.owner_id);
          }
        }
      );
    }
  }, [access_token, id]);

  if (!project) {
    return <div>Loading...</div>;
  }

  return (
    <div className="max-w-full overflow-x-clip">
      {
        //------------------------------------------
        // Project Basics Card
        //------------------------------------------
      }

      <div className="space-y-4 h-full">
        <div className="flex gap-2">
          <div>
            {/* <Image
              width={800}
              // height={600}
              alt="NextUI hero Image"
              className="p-2"
              src="https://picsum.photos/200/300"
              // src="https://www.liquidplanner.com/wp-content/uploads/2019/04/HiRes-17.jpg"
            /> */}
            <Image
              width={800}
              height={600}
              alt="Project Image"
              className="p-2"
              src={getRandomBackground()}
            />
          </div>

          <div className="p-6 max-w-2xl w-1/2">
            <div className={`${heading({ size: "lg" })}`}>{project.title}</div>
            <div className="h-full pl-2 pt-1 space-y-4 relative">
              <div className={`${subheading({ size: "lg" })} max-w-10`}>
                {project.description}
              </div>
              <div className="space-y-4 min-h-24">
                <div className="wrap-text break-words">
                  {project.introduction}
                </div>
                <div className="p-2">
                  <a
                    href={`profile?id=${project.owner_id}`}
                    className="flex items-center space-x-3 pl-4"
                  >
                    <Avatar
                      icon={<AvatarIcon />}
                      classNames={{
                        base: "bg-gradient-to-br from-[#00B4DB] to-[#0083B0]",
                        icon: "text-black/80",
                      }}
                      className="w-12 h-12"
                      alt={project.owner?.name}
                    />
                    <span className="text-base font-medium">
                      {project.owner?.name}
                    </span>
                  </a>
                </div>
                <div className="w-full">
                <div className="px-4 pt-0 flex items-center flex-wrap gap-2">
                  {project.tags.map((tag, index) => (
                    <Chip key={index} size="sm" radius="sm" variant="bordered">
                      {tag}
                    </Chip>
                  ))}
                </div>
              </div>
              </div>

            </div>
          </div>
        </div>
        {
          //------------------------------------------
          // Content projectSection
          //------------------------------------------
        }
        <div>
          {project.sections.map((card) => (
            <div
              className="flex w-full flex-col space-y-6"
              // key
            >
              <div className="space-y-2">
                <div className="w-full space-y-2 px-6 pb-6 pt-2">
                  <div className={subheading({ size: "lg" })}>{card.title}</div>
                  <div className="wrap-text break-words">{card.body}</div>
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
