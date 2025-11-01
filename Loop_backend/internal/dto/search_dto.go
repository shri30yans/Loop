package dto

// SearchToolsResult represents a search result with metadata
type SearchToolsResult struct {
	Score        float64  `json:"score"`
	ProjectID    string   `json:"project_id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
	Domains      []string `json:"domains"`
	Features     []string `json:"features"`
}

// SearchRequest represents a search request
type SearchToolsRequest struct {
	Query string `json:"query" validate:"required,min=2,max=500"`
	Mode  string `json:"mode" validate:"omitempty,oneof=vector graph hybrid"` // vector, graph, hybrid
	Limit int    `json:"limit" validate:"omitempty,min=1,max=100"`
}

// SearchResponse represents a search response
type SearchResponse struct {
	Results []SearchToolsResult `json:"results"`
	Total   int                 `json:"total"`
}
