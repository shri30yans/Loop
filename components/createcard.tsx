import { useState } from 'react';
//import { v4 as uuidv4 } from 'uuid'; // Assuming you're using UUID for unique IDs

// Define a type for the card
type Card = {
  id: string;
  title: string;
  body: string;
  tags?: string[];
};

// Initial cards can be generated or fetched from an API
const initialCards: Card[] = [
  { id: "title", title: "", body: "", tags: [] },
  { id: "introduction", title: "", body: "" },
  { id: "1", title: "", body: "" }, // Example of using UUID for unique IDs
  { id: "2", title: "", body: "" },
];

const PageComponent = () => {
  const [cards, setCards] = useState<Card[]>(initialCards);

  const addNewCard = () => {
    const newCard: Card = {
      id: "3", // Generate a unique ID for each new card
      title: "",
      body: "",
    };
    setCards([...cards, newCard]);
  };

  // Component rendering logic remains the same
};