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
// import { Selection } from "react";

export default function CreatePage() {

  const refresh_token = useAuthStore((state) => state.refresh_token);

  const type = [
    { key: "ai", label: "AI/ML" },
    { key: "web", label: "Web" },
    { key: "mobile", label: "Mobile" },
    { key: "devops", label: "DevOps" },
    { key: "startup", label: "Startup" },
    { key: "cloud", label: "Cloud" },
  ];

  const initialProjectSection: ProjectSectionType[] = [
    { section_number: 1, title: "", body: "" },
  ];

  const [projectSection, setProjectSection] = useState<ProjectSectionType[]>(initialProjectSection);

  const initialProject: ProjectType = {
    title: "",
    description: "",
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
    
    // Special handling for tags array
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
    const sectionNumber = parseInt(id, 10);
  
    const updatedSections = projectSection.map((section) => {
      if (section.section_number === sectionNumber) {
        return { ...section, [name]: value };
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
      //console.log("You can't add more than 10 updates")
      return;
    }
    const newCard: ProjectSectionType = {
      section_number: projectSection.length + 1,
      title: "",
      body: "",
    };
    setProjectSection([...projectSection, newCard]);
  };

  const handlePublish = (event: any) => {
    event.preventDefault();
    const user_id = useAuthStore.getState().user_id;
    if (user_id) {
      project.owner_id = user_id;
    }
    if (refresh_token){
    createProject(refresh_token,project);
  }
    
    // Reset all fields 
    setProject(initialProject);
    setProjectSection(initialProjectSection);
  };

  return (
    <div>
      <form className="space-y-4" onSubmit={handlePublish}>
        {
          //------------------------------------------
          // Project Basics Card
          //------------------------------------------
        }

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

                <div className="w-full space-y-2 px-6 ">
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
                    // maxLength={15}
                    // minLength={2}
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
                    // maxLength={30}
                    // minLength={10}
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
                      const tagsArray = Array.from(keys).map(String); // Convert numbers to strings
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
          {
            //------------------------------------------
            // Introduction Card
            //------------------------------------------
          }
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
                  // maxLength={250}
                  // minLength={50}
                  required
                />
              </div>
            </div>
          </CardBody>
        </Card>
        {
          //------------------------------------------
          // Content projectSection
          //------------------------------------------
        }
        {projectSection.map((card) => (
          <div className="flex w-full flex-col space-y-6" key={card.section_number}>
            <Card isBlurred>
              <CardBody>
                <div className="space-y-2">
                  <div className="w-full space-y-2 px-6 pb-6 pt-2">
                    <div className={subheading({ size: "lg" })}>
                      Update {card.section_number}
                    </div>
                    <Input
                      isRequired
                      className="w-full"
                      type="text"
                      label="Title"
                      id={card.section_number.toString()}
                      name="title"
                      value={card.title}
                      onChange={handleProjectSectionChange}
                      // maxLength={80}
                      // minLength={6}
                      required

                    />
                    <Textarea
                      label="Body"
                      className="w-full"
                      isRequired
                      id={card.section_number.toString()}
                      name="body"
                      value={card.body}
                      onChange={handleProjectSectionChange}
                      // maxLength={2000}
                      // minLength={50}
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
