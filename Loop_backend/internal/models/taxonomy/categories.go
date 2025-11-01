package taxonomy

// Category maps for standardized classification
var TechnologyCategories = map[string][]string{
    "Web Development": {"Frontend", "Backend", "Full Stack"},
    "Mobile Development": {"iOS", "Android", "Cross-Platform"},
    "Cloud Computing": {"AWS", "Azure", "GCP", "Serverless"},
    "Data Engineering": {"Databases", "Big Data", "ETL"},
    "AI/ML": {"Machine Learning", "Deep Learning", "NLP"},
    "DevOps": {"CI/CD", "Container Orchestration", "Infrastructure"},
}

var DomainCategories = map[string][]string{
    "Business": {"E-commerce", "Enterprise", "Finance", "Marketing"},
    "Social": {"Social Media", "Community", "Communication"},
    "Education": {"E-learning", "Training", "Assessment"},
    "Healthcare": {"Medical", "Wellness", "Healthcare IT"},
    "Entertainment": {"Gaming", "Media", "Content"},
}

var FeatureCategories = map[string][]string{
    "User Management": {"Authentication", "Authorization", "User Profiles"},
    "Data Management": {"Storage", "Analytics", "API"},
    "Communication": {"Real-time", "Messaging", "Notifications"},
    "Content": {"Content Management", "File Sharing", "Media Processing"},
    "Integration": {"Third-party APIs", "Payment Systems", "External Services"},
}

// CategoryTypes maps top-level category types to their subcategories
var CategoryTypes = map[string]map[string][]string{
    "Technology": TechnologyCategories,
    "Domain": DomainCategories,
    "Feature": FeatureCategories,
}

// Helper functions to check category validity
func IsValidTechnologyCategory(main, sub string) bool {
    if subCats, exists := TechnologyCategories[main]; exists {
        for _, validSub := range subCats {
            if validSub == sub {
                return true
            }
        }
    }
    return false
}

func IsValidDomainCategory(main, sub string) bool {
    if subCats, exists := DomainCategories[main]; exists {
        for _, validSub := range subCats {
            if validSub == sub {
                return true
            }
        }
    }
    return false
}

func IsValidFeatureCategory(main, sub string) bool {
    if subCats, exists := FeatureCategories[main]; exists {
        for _, validSub := range subCats {
            if validSub == sub {
                return true
            }
        }
    }
    return false
}

// IsValidCategory checks if a category is valid across all category types
func IsValidCategory(categoryType, main, sub string) bool {
    if categoryMap, exists := CategoryTypes[categoryType]; exists {
        if subCats, exists := categoryMap[main]; exists {
            for _, validSub := range subCats {
                if validSub == sub {
                    return true
                }
            }
        }
    }
    return false
}
