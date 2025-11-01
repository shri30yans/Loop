package models

import "time"

// Contribution represents a user's contribution to an entity
type Contribution struct {
    ActivityType  string      // Type of activity (CREATED, CONTRIBUTED_TO, etc.)
    TargetName    string      // Name of the target entity
    TargetType    string      // Type of the target entity
    Timestamp     time.Time   // When the contribution was made
    Properties    map[string]interface{} // Additional contribution properties
}
