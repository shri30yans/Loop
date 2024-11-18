package models

import (
	"time"
)

// Event represents an event in the database.
type Event struct {
	EventID int    `json:"event_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
}

// UserEventParticipation represents the participation of a user in an event.
type UserEventParticipation struct {
	UserID  int `json:"user_id"`
	EventID int `json:"event_id"`
}

type BaseUser struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    Location  *string   `json:"location"`
    Bio       *string   `json:"bio"`
}

type User struct {
    BaseUser
    HashedPassword string    `json:"hashed_password"`
    Projects       []Project `json:"projects"`
	Status 			string    `json:"status"`
}

type UserDetails struct {
    BaseUser
}

type UserInfoSummary struct {
    BaseUser
    Projects []ProjectsResponse `json:"projects"`
}
