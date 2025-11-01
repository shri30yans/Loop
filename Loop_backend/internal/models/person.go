package models

// Person represents a user or contributor in the graph database
type Person struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Email       string                 `json:"email"`
    Properties  map[string]interface{} `json:"properties"`
}
