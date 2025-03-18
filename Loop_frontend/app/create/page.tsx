"use client";
import { subheading } from "@/components/ui/primitives";
import { Textarea, Input } from "@nextui-org/input";
import { Tab, Tabs } from "@nextui-org/tabs";
import { Card, CardBody } from "@nextui-org/card";
import { Button } from "@nextui-org/button";
import { useState } from "react";
import { createProject } from "./actions";
import { Select, SelectItem } from "@nextui-org/select";
import { Image } from "@nextui-org/image";
import { ProjectSectionType, ProjectType } from "../types";
import { useAuthStore } from '../../lib/auth/authStore';
import { useRouter } from 'next/navigation';

export default function CreatePage() {
  const router = useRouter();
  const access_token = useAuthStore((state) => state.access_token);

  const type = [
    { key: "ai", label: "AI/ML" },
    { key: "web", label: "Web" },
    { key: "mobile", label: "Mobile" },
    { key: "devops", label: "DevOps" },
    { key: "startup", label: "Startup" },
    { key: "cloud", label: "Cloud" },
  ];

  const initialProjectSection: ProjectSectionType[] = [
    { index: 1, title: "", content: "" },
  ];

  const [projectSection, setProjectSection] = useState<ProjectSectionType[]>(initialProjectSection);

  const initialProject: ProjectType = {
    title: "",
    description: "",
    status: "draft",
    introduction: "",
    sections: projectSection,
    owner_id: "",
    tags: [],
  };
  
  const [project, setProject] = useState<ProjectType>(initialProject);
  
  const handleProjectChange = (
    event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement> | 
    { target: { name: string; value: string | string[] } }
  ) => {
    const { name, value } = event.target;
    
    if (name === 'tags') {
      const tagsArray = Array.isArray(value) ? value : value.split(',').map(tag => tag.trim());
      setProject(prevProject => ({
        ...prevProject,
        tags: tagsArray
      }));
    } else {
      setProject(prevProject => ({
        ...prevProject,
        [name]: value
      }));
    }
  };

  const handleProjectSectionChange = (event: any) => {
    const { name, value, id } = event.target;
    const sectionNumber = parseInt(id || '0', 10);
  
    const updatedSections = projectSection.map((section) => {
      if (section.index === sectionNumber) {
        const fieldName = name === 'body' ? 'content' : name;
        return { ...section, [fieldName]: value };
      }
      return section;
    });
  
    setProjectSection(updatedSections);
    setProject((prevProject) => ({
      ...prevProject,
      sections: updatedSections,
    }));
  };
  
  const addNewCard = () => {
    if (projectSection.length >= 10) {
      return;
    }
    const newCard: ProjectSectionType = {
      index: projectSection.length + 1,
      title: "",
      content: "",
    };
    setProjectSection([...projectSection, newCard]);
  };

  const handlePublish = async (event: any) => {
    event.preventDefault();
    const user_id = useAuthStore.getState().user_id;

    if (access_token && user_id) {
      try {
        project.owner_id = user_id;
        const sectionsWithoutIndex = project.sections.map(({ title, content }) => ({ title, content }));
        const response = await createProject(access_token, { 
          ...project, 
          sections: sectionsWithoutIndex as any,
          status: "draft"
        });

        // Debug the response
        console.log('Project creation response:', response);

        // Get ID from response (using just 'id' now since we've updated the types)
        const projectId = response.id;
        if (!projectId) {
          throw new Error('No project ID in response');
        }

        setProject(initialProject);
        setProjectSection(initialProjectSection);
        router.push(`/projectpage?id=${projectId}`);
      } catch (error) {
        console.error('Error during project creation:', error);
      }
    }
  };

  return (
    <div>
      <form className="space-y-4" onSubmit={handlePublish}>
        <Card isBlurred>
          <CardBody>
            <div className="space-y-2">
              <div className="flex items-center">
                <Image
                  width={600}
                  height={400}
                  alt="NextUI hero Image"
                  className="p-2"
                  src="https://www.liquidplanner.com/wp-content/uploads/2019/04/HiRes-17.jpg"
                />

                <div className="w-full space-y-2 px-6">
                  <div className={subheading({ size: "lg" })}>
                    Project basics
                  </div>
                  <Input
                    isRequired
                    className="w-full"
                    type="text"
                    label="Title"
                    name="title"
                    value={project.title}
                    onChange={handleProjectChange}
                    required
                  />
                  <Input
                    isRequired
                    className="w-full"
                    type="text"
                    label="Description"
                    name="description"
                    value={project.description}
                    onChange={handleProjectChange}
                    required
                  />

                  <Select
                    isRequired
                    label="Tags"
                    selectionMode="multiple"
                    className="w-full"
                    placeholder="What is your project about?"
                    name="tags"
                    onSelectionChange={(keys) => {
                      const tagsArray = Array.from(keys).map(String);
                      handleProjectChange({
                        target: {
                          name: 'tags',
                          value: tagsArray
                        }
                      });
                    }}
                    required
                  >
                    {type.map((data) => (
                      <SelectItem key={data.key} value={data.key}>
                        {data.label}
                      </SelectItem>
                    ))}
                  </Select>
                </div>
              </div>
            </div>
          </CardBody>
        </Card>

        <Card isBlurred>
          <CardBody>
            <div className="space-y-2">
              <div className="w-full space-y-2 px-6 pb-6 pt-2">
                <div className={subheading({ size: "lg" })}>
                  Introduction
                </div>
                <Textarea
                  label="Body"
                  className="w-full"
                  isRequired
                  placeholder="Explain why you made your project."
                  name="introduction"
                  value={project.introduction}
                  onChange={handleProjectChange}
                  required
                />
              </div>
            </div>
          </CardBody>
        </Card>

        {projectSection.map((card) => (
          <div className="flex w-full flex-col space-y-6" key={card.index}>
            <Card isBlurred>
              <CardBody>
                <div className="space-y-2">
                  <div className="w-full space-y-2 px-6 pb-6 pt-2">
                    <div className={subheading({ size: "lg" })}>
                      Update {card.index}
                    </div>
                    <Input
                      isRequired
                      className="w-full"
                      type="text"
                      label="Title"
                      id={card.index.toString()}
                      name="title"
                      value={card.title}
                      onChange={handleProjectSectionChange}
                      required
                    />
                    <Textarea
                      label="Body"
                      className="w-full"
                      isRequired
                      id={card.index.toString()}
                      name="body"
                      value={card.content}
                      onChange={handleProjectSectionChange}
                      required
                    />
                  </div>
                </div>
              </CardBody>
            </Card>
          </div>
        ))}
        <div className="flex gap-4">
          <Button
            className="w-3/4"
            color="primary"
            radius="lg"
            variant="flat"
            onPress={addNewCard}
          >
            Add New Step
          </Button>
          <Button
            type="submit"
            className="w-1/4"
            color="success"
            radius="lg"
            variant="flat"
          >
            Publish
          </Button>
        </div>
      </form>
    </div>
  );
}