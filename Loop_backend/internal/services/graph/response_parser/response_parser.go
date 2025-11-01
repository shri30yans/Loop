package response_parser

import (
    "Loop_backend/internal/models"
    "Loop_backend/internal/models/taxonomy"
    "fmt"
    "regexp"
    "strings"
)

var (
    // Pattern to match category_type:main_category:sub_category:justification format
    entityPattern = regexp.MustCompile(`\(([^:]+):([^:]+):([^:]+):([^)]+)\)`)
    keywordPattern = regexp.MustCompile(`"content_keywords":([^)]+)`)
)

// ResponseParser defines the interface for parsing responses into knowledge graphs
type ResponseParser interface {
    GenerateKnowledgeGraph(text string) (*models.KnowledgeGraph, error)
    ParseEntities(text string) ([]models.Entity, error)
    ExtractKeywords(text string) ([]string, error)
}

type responseParser struct{}

// NewResponseParser creates a new instance of ResponseParser
func NewResponseParser() ResponseParser {
    return &responseParser{}
}

func (p *responseParser) GenerateKnowledgeGraph(text string) (*models.KnowledgeGraph, error) {
    entities, err := p.ParseEntities(text)
    if err != nil {
        return nil, fmt.Errorf("error parsing entities: %w", err)
    }

    keywords, err := p.ExtractKeywords(text)
    if err != nil {
        return nil, fmt.Errorf("error extracting keywords: %w", err)
    }

    return &models.KnowledgeGraph{
        Entities: entities,
        Keywords: keywords,
    }, nil
}

func (p *responseParser) ParseEntities(text string) ([]models.Entity, error) {



    fmt.Println("=== START OF ORIGINAL INPUT ===")
    fmt.Println(text)
    fmt.Println("=== END OF ORIGINAL INPUT ===")

    lines := strings.Split(strings.TrimSpace(text), "\n")
    var entities []models.Entity
    
    fmt.Println("Found lines to process:", len(lines))

    // Track main categories to avoid duplicates
    mainCategories := make(map[string]bool)

    for i, line := range lines {
        trimmedLine := strings.TrimSpace(line)
        if trimmedLine == "" {
            continue
        }
        
        fmt.Println("\nProcessing line", i, ":", trimmedLine)
        
        if !strings.HasPrefix(trimmedLine, "(") || !strings.HasSuffix(trimmedLine, ")") {
            fmt.Println("Line", i, "doesn't match expected format (skipping):", trimmedLine)
            continue
        }
        
        matches := entityPattern.FindAllStringSubmatch(trimmedLine, -1)
        fmt.Println("Line", i, "matches:", matches)
        
        for _, match := range matches {
            if len(match) != 5 { // Updated to expect 5 parts (full match + 4 capture groups)
                continue
            }

            categoryType := strings.TrimSpace(match[1])
            mainCategory := strings.TrimSpace(match[2])
            subCategory := strings.TrimSpace(match[3])
            justification := strings.TrimSpace(match[4])

            // Validate the category using the new validation function
            if !taxonomy.IsValidCategory(categoryType, mainCategory, subCategory) {
                fmt.Printf("Invalid category combination: type=%s, main=%s, sub=%s\n", 
                    categoryType, mainCategory, subCategory)
                continue
            }

            // Create main category entity if it doesn't exist
            if !mainCategories[mainCategory] {
                mainEntity := models.Entity{
                    Name:        mainCategory,
                    Type:        categoryType,
                    Level:       models.LevelMain,
                    Description: fmt.Sprintf("Main category for %s", categoryType),
                }
                entities = append(entities, mainEntity)
                mainCategories[mainCategory] = true
            }

            // Create sub-category entity
            subEntity := models.Entity{
                Name:         subCategory,
                Type:        categoryType,
                Level:       models.LevelSub,
                ParentEntity: mainCategory,
                Description: justification,
            }
            entities = append(entities, subEntity)
        }
    }

    return entities, nil
}

func (p *responseParser) ExtractKeywords(text string) ([]string, error) {
    match := keywordPattern.FindStringSubmatch(text)
    if len(match) != 2 {
        return nil, nil
    }

    keywordsStr := strings.TrimSpace(match[1])
    keywords := strings.Split(keywordsStr, ",")

    for i, keyword := range keywords {
        keywords[i] = strings.TrimSpace(keyword)
    }

    return keywords, nil
}

func ParseWeight(weightStr string) int {
    var weight int
    _, err := fmt.Sscanf(weightStr, "%d", &weight)
    if err != nil || weight < 1 || weight > 10 {
        weight = 5
    }
    return weight
}

