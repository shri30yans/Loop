"use client";

import { Button } from "@nextui-org/button";
import { Divider } from "@nextui-org/divider";

interface PostCardProps {
  title: string;
  body: string;
}

export default function PostCard({ title, body}: PostCardProps) {
  return (
    <div className="max-w-full mt-4 radius-full">
      <div>
        <div className="text-3xl font-semibold">{title}</div>
      </div>
      <div className="text-xl">{body}</div>
      <div className="flex flex-row mt-3 gap-2">
        <Button size="sm" variant="flat">
          Up
        </Button>
        <Button size="sm" variant="flat">
          Down
        </Button>
        <Button size="sm" variant="flat">
          Share
        </Button>
      </div>
      <Divider className="mt-4" />
    </div>
  );
}
