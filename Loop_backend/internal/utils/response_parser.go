package utils

import (
	"Loop_backend/internal/models"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

var (
	instance *ResponseParser
	once     sync.Once

	// Regex patterns
	entityPattern   = regexp.MustCompile(`\("entity":([^:]+):([^:]+):([^)]+)\)`)
	relationPattern = regexp.MustCompile(`\("relationship":([^:]+):([^:]+):([^:]+):([^:]+):([^:]+):([^)]+)\)`)
	keywordPattern  = regexp.MustCompile(`"content_keywords":([^)]+)`)
)

type ResponseParser struct{}

// GetResponseParser returns a singleton instance of ResponseParser
func GetResponseParser() *ResponseParser {
	once.Do(func() {
		instance = &ResponseParser{}
	})
	return instance
}

func (p *ResponseParser) GenerateKnowledgeGraph(text string) (*models.KnowledgeGraph, error) {
	entities, err := p.ParseEntities(text)
	if err != nil {
		return nil, fmt.Errorf("error parsing entities: %w", err)
	}

	relationships, err := p.ParseRelationships(text)
	if err != nil {
		return nil, fmt.Errorf("error parsing relationships: %w", err)
	}

	keywords, err := p.ExtractKeywords(text)
	if err != nil {
		return nil, fmt.Errorf("error extracting keywords: %w", err)
	}

	return &models.KnowledgeGraph{
		Entities:      entities,
		Relationships: relationships,
		Keywords:      keywords,
	}, nil
}

func (p *ResponseParser) ParseEntities(text string) ([]models.Entity, error) {
	matches := entityPattern.FindAllStringSubmatch(text, -1)
	entities := make([]models.Entity, 0, len(matches))

	for _, match := range matches {
		if len(match) != 4 {
			continue
		}

		name := strings.TrimSpace(match[1])
		entityType := strings.TrimSpace(match[2])
		desc := strings.TrimSpace(match[3])

		entity := models.Entity{
			Name:        name,
			Type:        entityType,
			Description: desc,
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

func (p *ResponseParser) ParseRelationships(text string) ([]models.Relationship, error) {
	matches := relationPattern.FindAllStringSubmatch(text, -1)
	relationships := make([]models.Relationship, 0, len(matches))

	for _, match := range matches {
		if len(match) != 7 {
			continue
		}

		source := strings.TrimSpace(match[1])
		target := strings.TrimSpace(match[2])
		desc := strings.TrimSpace(match[3])
		relType := strings.TrimSpace(match[4])
		weight := ParseWeight(match[5])
		category := strings.TrimSpace(match[6])

		// Format relationship type to be Neo4j compatible
		formattedRelType := strings.ToUpper(strings.ReplaceAll(relType, " ", "_"))

		relationship := models.Relationship{
			Source:      source,
			Target:      target,
			Description: desc,
			Type:        formattedRelType,
			Weight:      weight,
			Category:    category,
		}
		relationships = append(relationships, relationship)
	}

	return relationships, nil
}

func (p *ResponseParser) ExtractKeywords(text string) ([]string, error) {
	match := keywordPattern.FindStringSubmatch(text)
	if len(match) != 2 {
		return nil, nil
	}

	keywordsStr := strings.TrimSpace(match[1])
	keywords := strings.Split(keywordsStr, ",")

	// Clean and trim keywords
	for i, keyword := range keywords {
		keywords[i] = strings.TrimSpace(keyword)
	}

	return keywords, nil
}

func ParseWeight(weightStr string) int {
	// Convert string weight to int, default to 5 if invalid
	var weight int
	_, err := fmt.Sscanf(weightStr, "%d", &weight)
	if err != nil || weight < 1 || weight > 10 {
		weight = 5
	}
	return weight
}
