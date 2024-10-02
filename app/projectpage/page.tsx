"use client";
import { Card, CardBody } from "@nextui-org/card";
import { Image } from "@nextui-org/image";
import { subheading, heading } from "@/components/primitives";
import { ProjectType } from "../types";
import { useEffect, useState } from "react";
import { fetchProjectInfo } from "./actions";
import { Chip } from "@nextui-org/chip";

export default function ProjectPage() {
  const [project, setProject] = useState<ProjectType | null>(null);
  const queryParams = new URLSearchParams(window.location.search);
  const id = queryParams.get("id");
  useEffect(() => {
    // Ensure 'id' is not undefined or empty
    if (id) {
      fetchProjectInfo(String(id)).then(
        (fetchedProject: ProjectType | null) => {
          if (fetchedProject) {
            setProject(fetchedProject);
            console.log(fetchedProject);
          }
        }
      );
    }
  }, [id]); // Add 'id' as a dependency to useEffect

  if (!project) {
    return <div>Error!</div>;
  }

  return (
    // <div className="max-w-full overflow-x-clip">
    <div className="max-w-full overflow-x-clip">

      {
        //------------------------------------------
        // Project Basics Card
        //------------------------------------------
      }

      <div className="space-y-4">
        <div className="flex gap-2">
          <div>
          <Image
            width={800}
            height={600}
            alt="NextUI hero Image"
            className="p-2"
            src="https://nextui-docs-v2.vercel.app/images/hero-card-complete.jpeg"
          />
          </div>

          <div className="p-6 max-w-2xl">
            <div className={`${heading({ size: "lg" })}`}>{project.title}</div>
            <div className="h-full pl-4 pt-1 space-y-4 relative">
              <div className={`${subheading({ size: "lg" })} max-w-10`}>
                {project.description}
              </div>
              <div className=" wrap-text break-words">{project.introduction}</div>
              <div className="absolute bottom-10 p-2" >
                <Chip size="md" radius="sm" variant="bordered">
                  {project.tags}
                </Chip>
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
          <div className="flex w-full flex-col space-y-6" key={card.id}>
            <div className="space-y-2">
              <div className="w-full space-y-2 px-6 pb-6 pt-2">
                <div className={subheading({ size: "lg" })}>{card.title}</div>
                <div className="wrap-text break-words">{card.body}</div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
    </div>

  );
}
