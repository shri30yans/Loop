package response_parser

import (
	"Loop_backend/internal/models"
	"fmt"
	"regexp"
	"strings"
)

var (
	// Regex patterns
	entityPattern = regexp.MustCompile(`\(([^:]+):([^:]+):([^)]+)\)`)
	keywordPattern  = regexp.MustCompile(`"content_keywords":([^)]+)`)
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

func (p *responseParser) ExtractKeywords(text string) ([]string, error) {
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
