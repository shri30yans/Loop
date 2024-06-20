"use client";
import { title } from "@/components/primitives";
import { Textarea, Input } from "@nextui-org/input";
import { Tab, Tabs } from "@nextui-org/tabs";
import { Card, CardBody } from "@nextui-org/card";
import { Button } from "@nextui-org/button";
import { useState } from "react";
import { addPost } from "./actions";
import { Select, SelectItem } from "@nextui-org/select";

export default function CreatePage() {
  const type = [
    { key: "ai", label: "AI/ML" },
    { key: "web", label: "Web" },
    { key: "mobile", label: "Mobile" },
    { key: "devops", label: "DevOps" },
    { key: "startup", label: "Startup" },
    { key: "cloud", label: "Cloud" },
  ];

  const [post, setPost] = useState({ title: "", body: "" });

  const handleInputChange = (event: any) => {
    const { name, value } = event.target;
    setPost((prevState) => ({ ...prevState, [name]: value }));
  };

  const handleSubmit = (event: any) => {
    event.preventDefault();
    addPost(post);
    setPost({ title: "", body: "" });
  };

  return (
    <div className="flex w-full flex-col space-y-4">
      <h1 className={title()}>Create</h1>
      <Tabs aria-label="Options">
        <Tab key="posts" title="Posts">
          <Card>
            <CardBody>
              <form onSubmit={handleSubmit}>
                <div className="space-y-2">
                  <Input
                    isRequired
                    className="w-1/2"
                    type="text"
                    label="Title"
                    name="title"
                    value={post.title}
                    onChange={handleInputChange}
                  />
                  <Select
                    label="Tags"
                    selectionMode="multiple"
                    className="w-3/4"
                  >
                    {type.map((data) => (
                      <SelectItem key={data.key}>{data.label}</SelectItem>
                    ))}
                  </Select>
                  <Textarea
                    isRequired
                    label="Body"
                    className="w-full"
                    name="body"
                    value={post.body}
                    onChange={handleInputChange}
                  />
                  <Button
                    type="submit"
                    className="w-full"
                    color="primary"
                    radius="lg"
                    variant="flat"
                  >
                    {" "}
                    Submit{" "}
                  </Button>
                </div>
              </form>
            </CardBody>
          </Card>
        </Tab>

        <Tab key="projects" title="Projects">
          <Card>
            <CardBody>
              <div className="space-y-2">
                <Input isRequired className="w-1/2" type="text" label="Title" />
                <Input
                  isRequired
                  className="w-1/2"
                  type="text"
                  label="Describe your project in one line"
                />
                <Input className="w-1/2" type="text" label="Video Link" />
                <Input className="w-1/2" type="text" label="Image" />
                <Select label="Tags" selectionMode="multiple" className="w-3/4">
                  {type.map((data) => (
                    <SelectItem key={data.key}>{data.label}</SelectItem>
                  ))}
                </Select>
                <div>
                  <Textarea
                    label="Body"
                    className="w-full"
                    value={post.body}
                  />
                </div>
                <div>
                  <Button
                    className="w-full"
                    color="primary"
                    radius="lg"
                    variant="flat"
                  >
                    {" "}
                    Submit{" "}
                  </Button>
                </div>
              </div>
            </CardBody>
          </Card>
        </Tab>
      </Tabs>
    </div>
  );
}
