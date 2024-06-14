"use client";
import { title } from "@/components/primitives";
import { Textarea, Input } from "@nextui-org/input";
import { Tab, Tabs } from "@nextui-org/tabs";
import { Card, CardBody } from "@nextui-org/card";
import {Button} from "@nextui-org/button";
import { useState } from "react";
import {addPost} from "./actions";

export default function CreatePage() {
  const [postTitle, setPostTitle] = useState('');
  const [postBody, setPostBody] = useState('');
  const handleSubmit = (event :any) => {
    event.preventDefault();
    addPost({title: postTitle, body: postBody});
    setPostTitle('');
    setPostBody('');
    
  };
  return (
    <div className="flex w-full flex-col space-y-4">
      <h1 className={title()}>Create</h1>
      <Tabs aria-label="Options">
        <Tab key="posts" title="Posts">
          <Card>
            <CardBody>
              <form onSubmit={handleSubmit}>
                <div className="space-y-4">
                  <div>
                    <Input isRequired className="w-1/2" type="text" label="Title" value={postTitle} onChange={(e) => setPostTitle(e.target.value)} />
                  </div>
                  <div>
                    <Textarea isRequired label="Body" className="w-full" value={postBody} onChange={(e) => setPostBody(e.target.value)} />
                  </div>
                  <div>
                    <Button type="submit" className="w-full" color="primary" radius="lg" variant="flat"> Submit </Button>
                  </div>
                </div>
              </form>
            </CardBody>
          </Card>
        </Tab>

        <Tab key="projects" title="Projects">
          <Card>
            <CardBody>
              <div className="space-y-4">
                <div>
                  <Input className="w-1/2" type="text" label="Title" />
                </div>
                <div>
                  <Textarea label="Body" className="w-full" />
                </div>
                <div>
                  <Button className="w-full" color="primary" radius="lg" variant="flat"> Submit </Button>
                </div>
              </div>
              
            </CardBody>
          </Card>
        </Tab>

        <Tab key="blogs" title="Blogs">
          <Card>
            <CardBody>
              <div className="space-y-4">
                <div>
                  <Input className="w-1/2" type="text" label="Title" />
                </div>
                <div>
                  <Textarea label="Body" className="w-full" />
                </div>
                <div>
                  <Button className="w-full" color="primary" radius="lg" variant="flat"> Submit </Button>
                </div>
              </div>
            </CardBody>
          </Card>
        </Tab>
      </Tabs>
    </div>
  );
}
