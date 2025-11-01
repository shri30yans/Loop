package project

import (
    "Loop_backend/internal/models"
    "Loop_backend/internal/models/taxonomy"
    "fmt"
    "strings"
)

const (
    // Delimiters used in prompts and response parsing
    TupleDelimiter      = ":"       // Separates fields within a record
    RecordDelimiter     = "\n---\n" // Separates individual records in the response
    CompletionDelimiter = "[END]"   // Marks the end of the LLM response
    SummaryStartTag     = "[SUMMARY_START]"
    SummaryEndTag       = "[SUMMARY_END]"
)

// formatCategoryInstructions formats the category hierarchies for the prompt
func formatCategoryInstructions() string {
    var sb strings.Builder
    
    // Format Technology Categories
    sb.WriteString("\nTechnology Categories:\n")
    for mainCat, subCats := range taxonomy.TechnologyCategories {
        sb.WriteString(fmt.Sprintf("- %s:\n", mainCat))
        for _, subCat := range subCats {
            sb.WriteString(fmt.Sprintf("  * %s\n", subCat))
        }
    }
    
    // Format Domain Categories
    sb.WriteString("\nDomain Categories:\n")
    for mainCat, subCats := range taxonomy.DomainCategories {
        sb.WriteString(fmt.Sprintf("- %s:\n", mainCat))
        for _, subCat := range subCats {
            sb.WriteString(fmt.Sprintf("  * %s\n", subCat))
        }
    }
    
    // Format Feature Categories
    sb.WriteString("\nFeature Categories:\n")
    for mainCat, subCats := range taxonomy.FeatureCategories {
        sb.WriteString(fmt.Sprintf("- %s:\n", mainCat))
        for _, subCat := range subCats {
            sb.WriteString(fmt.Sprintf("  * %s\n", subCat))
        }
    }
    
    return sb.String()
}

// GetProjectAnalysisPrompt returns the prompt for project analysis
func GetProjectAnalysisPrompt(project *models.Project) string {
    // Get formatted categories
    categories := formatCategoryInstructions()
    
    // Combine project sections
    var combinedSections strings.Builder
    for _, section := range project.Sections {
        combinedSections.WriteString(section.Content)
        combinedSections.WriteString("\n")
    }
    text := combinedSections.String()

    return fmt.Sprintf(`---Goal---
Analyze the project and categorize it according to the predefined categories below.

---Available Categories---%s

---Instructions---
1. For each category type (Technology, Domain, Feature), identify the most relevant 
   categories that apply to this project. Format each category selection as:
   (<main_category>%s<sub_category>%s<brief_justification>)

   Example format:
   (Web Development:Frontend:Project extensively uses React and modern frontend technologies)
   (Business:E-commerce:Implements online shopping and payment processing features)

   Guidelines:
   - Choose ONLY from the provided categories and subcategories
   - Select the most specific and relevant categories
   - One category per line
   - No spaces around delimiters
   - Include a clear justification for each selection
   - You must select at least one Technology category and one Domain category
   - Maximum of 3 categories per type (Technology, Domain, Feature)

2. After categorization, provide a concise project summary:
   %s
   <project summary emphasizing the selected categories and their relationships>
   %s

Only provide categories and summary in the specified formats. Do not include any additional explanation.

---Input Text---
%s`,
        categories,
        TupleDelimiter,
        TupleDelimiter,
        SummaryStartTag,
        SummaryEndTag,
        text,
    )
}
