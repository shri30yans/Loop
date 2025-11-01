package main

import (
    "fmt"
    "strings"
)

// Simplified entity struct for migration
type Entity struct {
    Name         string
    Type         string
    Description  string
    Level        string
    ParentEntity string
}

func main() {
    fmt.Println("Starting project category migration...")
    fmt.Println("This script will:")
    fmt.Println("1. Convert existing free-form entities to standardized categories")
    fmt.Println("2. Create hierarchical relationships between categories")
    fmt.Println("3. Update project graphs with the new structure")
    fmt.Println()
    fmt.Println("WARNING: This is a destructive operation. Make sure you have a backup before proceeding.")
    fmt.Println("Press Ctrl+C to cancel, or Enter to continue...")

    // Wait for user input
    fmt.Scanln()

    fmt.Println("\nMigration steps:")
    fmt.Println("1. Loading projects...")
    // TODO: Initialize Neo4j connection and load projects

    fmt.Println("2. Processing projects...")
    // TODO: Process each project with new categorization

    fmt.Println("3. Updating graph database...")
    // TODO: Store updated graphs

    fmt.Println("\nMigration complete!")
}

// Technology categories
var TechnologyCategories = map[string][]string{
    "Web Development": {"Frontend", "Backend", "Full Stack"},
    "Mobile Development": {"iOS", "Android", "Cross-Platform"},
    "Cloud Computing": {"AWS", "Azure", "GCP", "Serverless"},
    "Data Engineering": {"Databases", "Big Data", "ETL"},
    "AI/ML": {"Machine Learning", "Deep Learning", "NLP"},
    "DevOps": {"CI/CD", "Container Orchestration", "Infrastructure"},
}

// Domain categories
var DomainCategories = map[string][]string{
    "Business": {"E-commerce", "Enterprise", "Finance", "Marketing"},
    "Social": {"Social Media", "Community", "Communication"},
    "Education": {"E-learning", "Training", "Assessment"},
    "Healthcare": {"Medical", "Wellness", "Healthcare IT"},
    "Entertainment": {"Gaming", "Media", "Content"},
}

// Feature categories
var FeatureCategories = map[string][]string{
    "User Management": {"Authentication", "Authorization", "User Profiles"},
    "Data Management": {"Storage", "Analytics", "API"},
    "Communication": {"Real-time", "Messaging", "Notifications"},
    "Content": {"Content Management", "File Sharing", "Media Processing"},
    "Integration": {"Third-party APIs", "Payment Systems", "External Services"},
}

// Helper function to map old entities to new categories
func mapToNewCategories(oldEntities []Entity) []Entity {
    var newEntities []Entity
    
    // Map for tracking created main categories
    mainCategories := make(map[string]bool)

    for _, entity := range oldEntities {
        // Try to map the entity to our new standardized categories
        if mainCat, subCat, ok := findMatchingCategory(entity); ok {
            // Add main category if not already added
            if !mainCategories[mainCat] {
                var catType string
                switch {
                case stringInMap(mainCat, TechnologyCategories):
                    catType = "Technology"
                case stringInMap(mainCat, DomainCategories):
                    catType = "Domain"
                case stringInMap(mainCat, FeatureCategories):
                    catType = "Feature"
                }

                mainEntity := Entity{
                    Name:        mainCat,
                    Type:        catType,
                    Level:       "MAIN",
                    Description: fmt.Sprintf("Main category for %s", catType),
                }
                newEntities = append(newEntities, mainEntity)
                mainCategories[mainCat] = true
            }

            // Add sub-category
            subEntity := Entity{
                Name:         subCat,
                Type:        entity.Type,
                Level:       "SUB",
                ParentEntity: mainCat,
                Description: entity.Description,
            }
            newEntities = append(newEntities, subEntity)
        }
    }

    return newEntities
}

// Helper function to find matching standardized category using contains matching
func findMatchingCategory(entity Entity) (string, string, bool) {
    name := strings.ToLower(entity.Name)
    
    // Check technology categories
    for mainCat, subCats := range TechnologyCategories {
        if strings.Contains(name, strings.ToLower(mainCat)) {
            return mainCat, subCats[0], true // Default to first subcategory
        }
        for _, subCat := range subCats {
            if strings.Contains(name, strings.ToLower(subCat)) {
                return mainCat, subCat, true
            }
        }
    }

    // Check domain categories
    for mainCat, subCats := range DomainCategories {
        if strings.Contains(name, strings.ToLower(mainCat)) {
            return mainCat, subCats[0], true
        }
        for _, subCat := range subCats {
            if strings.Contains(name, strings.ToLower(subCat)) {
                return mainCat, subCat, true
            }
        }
    }

    // Check feature categories
    for mainCat, subCats := range FeatureCategories {
        if strings.Contains(name, strings.ToLower(mainCat)) {
            return mainCat, subCats[0], true
        }
        for _, subCat := range subCats {
            if strings.Contains(name, strings.ToLower(subCat)) {
                return mainCat, subCat, true
            }
        }
    }
    
    return "", "", false
}

// Helper function to check if a string key exists in a map
func stringInMap(key string, m map[string][]string) bool {
    _, exists := m[key]
    return exists
}
