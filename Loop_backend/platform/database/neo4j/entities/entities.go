package entities

import (
    "strings"
)

// EntityTypeInfo contains metadata about an entity type
type EntityTypeInfo struct {
    Name string
    Description string
    ExtractionGuidelines string
}

// Entity Types
const (
    TypeProject     = "Project"
    TypeTechnology  = "Technology"
    TypeFeature    = "Feature"
    TypeTag        = "Tag"
    TypeCategory   = "Category"
    TypeStakeholder = "Stakeholder"
    TypePerson     = "Person"
    TypePlatform   = "Platform"
    TypeMethodology = "Methodology"
)

// EntityTypeMetadata provides descriptions and guidelines for each entity type
var EntityTypeMetadata = map[string]EntityTypeInfo{
    TypeProject: {
        Name: "Project",
        Description: "Main project or system being described",
        ExtractionGuidelines: "Look for project names, system titles, or main application identifiers",
    },
    TypeTechnology: {
        Name: "Technology",
        Description: "Technical tools, languages, frameworks used in the project",
        ExtractionGuidelines: "Identify programming languages, databases, libraries, frameworks, or tools",
    },
    TypeFeature: {
        Name: "Feature",
        Description: "Functionality or capabilities of the project",
        ExtractionGuidelines: "Look for described functionalities, capabilities, or user-facing features",
    },
    TypeTag: {
        Name: "Tag",
        Description: "Keywords or labels that categorize the project",
        ExtractionGuidelines: "Extract keywords, themes, or categorical terms",
    },
    TypeCategory: {
        Name: "Category",
        Description: "High-level classification or grouping",
        ExtractionGuidelines: "Identify broad categories or domains the project belongs to",
    },
    TypeStakeholder: {
        Name: "Stakeholder",
        Description: "Business roles or departments involved",
        ExtractionGuidelines: "Look for business roles, departments, or organizational units",
    },
    TypePerson: {
        Name: "Person",
        Description: "Individual contributors or team members",
        ExtractionGuidelines: "Extract names of people, developers, or team members",
    },
    TypePlatform: {
        Name: "Platform",
        Description: "Operating environment or system platform",
        ExtractionGuidelines: "Identify operating systems, cloud platforms, or deployment environments",
    },
    TypeMethodology: {
        Name: "Methodology",
        Description: "Development approaches or methods used",
        ExtractionGuidelines: "Look for development methodologies, processes, or practices",
    },
}

// EntityTypes returns a slice of all valid entity types
func EntityTypes() []string {
    types := make([]string, 0, len(EntityTypeMetadata))
    for typeName := range EntityTypeMetadata {
        types = append(types, typeName)
    }
    return types
}

// IsValidEntityType checks if the given type is a valid entity type
func IsValidEntityType(entityType string) bool {
    for _, validType := range EntityTypes() {
        if strings.EqualFold(validType, entityType) {
            return true
        }
    }
    return false
}

// GetEntityTypeInfo returns the metadata for a given entity type
func GetEntityTypeInfo(entityType string) (EntityTypeInfo, bool) {
    info, exists := EntityTypeMetadata[entityType]
    return info, exists
}
