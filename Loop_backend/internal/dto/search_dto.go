// filepath: d:\sem_6\code\genai-main-1\Loop\Loop_backend\internal\dto\search_dto.go
package dto

type SearchRequest struct {
	Query string `json:"query" validate:"required"`
}
