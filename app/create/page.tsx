"use client";
import { title, sectionheading } from "@/components/primitives";
import { Textarea, Input } from "@nextui-org/input";
import { Tab, Tabs } from "@nextui-org/tabs";
import { Card, CardBody } from "@nextui-org/card";
import { Button } from "@nextui-org/button";
import { useState } from "react";
import { addPost } from "./actions";
import { Select, SelectItem } from "@nextui-org/select";
import { Image } from "@nextui-org/image";
import { v4 as uuidv4 } from 'uuid';

export default function CreatePage() {
  const type = [
    { key: "ai", label: "AI/ML" },
    { key: "web", label: "Web" },
    { key: "mobile", label: "Mobile" },
    { key: "devops", label: "DevOps" },
    { key: "startup", label: "Startup" },
    { key: "cloud", label: "Cloud" },
  ];

  type Card = {
    id: number;
    title: string;
    body: string;
    tags?: string[];
  };


  const handleInputChange = (event: any) => {
    const { name, value, id } = event.target;
    console.log(name, value);
    setCards(cards => cards.map((card) => {
      if ( card.id === id) { 
        console.log("3e3",card);
        return { ...card, ['title']: value }; // Correctly update the field based on `name`
      }
      return card;
    }));
    console.log(name,value,cards);
  };

// Initial cards can be generated or fetched from an API
const initialCards: Card[] = [
  { id: 0, title: "", body: "", tags: [] },
  { id: 1, title: "", body: "" },
  { id: 2, title: "", body: "" }, // Example of using UUID for unique IDs
];
const [cards, setCards] = useState<Card[]>(initialCards);

const addNewCard = () => {
  if (cards.length >= 10) {
    console.log("You can't add more than 10 updates")
    return;
  }
  const newCard: Card = {
    id: cards.length, // Generate a unique ID for each new card
    title: "",
    body: "",
  };
  setCards([...cards, newCard]);
};

  return (
    <div className="space-y-4">
      <Card>
        <CardBody>
          <div className="space-y-2">
            <div className="flex items-center">
              <Image
                width={600}
                height={400}
                alt="NextUI hero Image"
                className="p-2"
                src="https://nextui-docs-v2.vercel.app/images/hero-card-complete.jpeg"
              />

              <div className="w-full space-y-2 px-6 ">
                <div className={sectionheading({ size: "lg" })}>
                  Project basics
                </div>
                <Input
                  isRequired
                  className="w-full"
                  type="text"
                  label="Title"
                  // value={tite}l
                  id="0"
                  name="title"
                  onChange={handleInputChange}
                />
                <Input
                  isRequired
                  className="w-full"
                  type="text"
                  label="Description"
                  name="title"
                  id="description"
                  value="description"
                  onChange={handleInputChange}
                />

                <Select
                  label="Tags"
                  selectionMode="single"
                  className="w-full"
                  placeholder="What is your project about?"
                >
                  {type.map((data) => (
                    <SelectItem key={data.key}>{data.label}</SelectItem>
                  ))}
                </Select>
              </div>
            </div>
          </div>
        </CardBody>
      </Card>
      <Card>
        <CardBody>
          <div className="space-y-2">
            <div className="w-full space-y-2 px-6 pb-6 pt-2">
              <div className={sectionheading({ size: "lg" })}>
               Introduction
              </div>
              <Input isRequired className="w-full" type="text" label="Title" />
              <Textarea label="Body" className="w-full" isRequired placeholder="Explain why you made your project."/>
            </div>
          </div>
        </CardBody>
      </Card>
      {cards.slice(1).map((card) => (
        <div className="flex w-full flex-col space-y-6 ">
          <Card key={card.id}>
            <CardBody>
              <div className="space-y-2">
                <div className="w-full space-y-2 px-6 pb-6 pt-2">
                  <div className={sectionheading({ size: "lg" })}>
                    Update {card.id}
                  </div>
                  <Input
                    isRequired
                    className="w-full"
                    type="text"
                    label="Title"
                  />
                  <Textarea label="Body" className="w-full" isRequired/>
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
        onClick={addNewCard}
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
    </div>
  );
}
