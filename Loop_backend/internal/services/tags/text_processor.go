package services

import (
    "strings"
    "unicode"
)

type TextProcessor interface {
    ExtractKeywords(text string) []string
    NormalizeText(text string) string
}

type textProcessor struct{}

func NewTextProcessor() TextProcessor {
    return &textProcessor{}
}

func (p *textProcessor) ExtractKeywords(text string) []string {
    normalizedText := p.NormalizeText(text)
    words := strings.Fields(normalizedText)

    // Simple keyword extraction: filter out short words
    var keywords []string
    for _, word := range words {
        if len(word) > 3 { // Example: only include words longer than 3 characters
            keywords = append(keywords, word)
        }
    }
    return keywords
}

func (p *textProcessor) NormalizeText(text string) string {
    // Convert to lowercase and remove non-alphanumeric characters
    var builder strings.Builder
    for _, r := range text {
        if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
            builder.WriteRune(unicode.ToLower(r))
        }
    }
    return builder.String()
}
