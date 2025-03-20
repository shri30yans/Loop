"use client";

import { Button } from "@nextui-org/button";
import { Divider } from "@nextui-org/divider";
import { Card, CardFooter, CardBody } from "@nextui-org/card";
import { Image } from "@nextui-org/image";
import { Chip } from "@nextui-org/chip";
import { Skeleton } from "@nextui-org/skeleton";

interface ProjectCardProps {
  isLoaded: boolean;
  title: string;
  body: string;
  tags?: string[];
}

export default function ProjectCard({
  title,
  body,
  tags = [],
}: ProjectCardProps) {
  return (
    <div>
      <Card isPressable isBlurred className="my-4">
        <CardBody className="overflow-visible ">
          <div>
            <Image
              alt="Card background"
              className="object-cover rounded-xl pb-2"
              src="https://www.liquidplanner.com/wp-content/uploads/2019/04/HiRes-17.jpg"
            />
          </div>
          <div className="text-2xl font-semibold">{title}</div>
          <div className="text-sm">{body}</div>
        </CardBody>
        <CardFooter className="px-4 pt-0 flex items-center flex-wrap gap-2">
          {tags.map((tag, index) => (
            <Chip key={index} size="sm" radius="sm" variant="bordered">
              {tag}
            </Chip>
          ))}
        </CardFooter>
      </Card>
    </div>
  );
}
