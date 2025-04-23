package handlers

import (
	"Loop_backend/internal/ai/interfaces"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ChatHandler handles chat-related requests
type ChatHandler struct {
	searchService  services.QueryService
	summaryService services.SummaryService
	provider       interfaces.Provider // Your LLM provider
}

// NewChatHandler creates a new chat handler
func NewChatHandler(
	searchService services.QueryService,
	summaryService services.SummaryService,
	provider interfaces.Provider,
) *ChatHandler {
	return &ChatHandler{
		searchService:  searchService,
		summaryService: summaryService,
		provider:       provider,
	}
}

// RegisterRoutes registers all routes for the chat handler
func (h *ChatHandler) RegisterRoutes(r RouteRegister) {
	r.RegisterProtectedRoute("/api/chat/conversation", h.HandleConversation, nil)
	r.RegisterProtectedRoute("/api/chat/history", h.GetChatHistory, nil)
}

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
}

// HandleConversation processes chat messages and routes them to the appropriate service
func (h *ChatHandler) HandleConversation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[ChatHandler] HandleConversation: Request received")

	// Parse request
	var chatReq ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&chatReq); err != nil {
		fmt.Printf("[ChatHandler] ERROR: Failed to decode request: %v\n", err)
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	fmt.Printf("[ChatHandler] Message received: %s\n", chatReq.Message)

	// Analyze the message to determine intent
	message := chatReq.Message
	messageContent := strings.ToLower(message)

	fmt.Println("[ChatHandler] Analyzing message intent...")

	// Create response object
	chatResp := ChatResponse{
		ID:        uuid.New().String(),
		Type:      "llm",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	fmt.Printf("[ChatHandler] Provider type: %T\n", h.provider)

	// Check for project summary intent
	if matchesSummaryIntent(messageContent) {
		fmt.Println("[ChatHandler] INTENT DETECTED: Summary request")

		projectName := extractProjectName(messageContent)
		fmt.Printf("[ChatHandler] Extracted project name: %s\n", projectName)

		if projectName != "" {
			// Search for the project first
			fmt.Printf("[ChatHandler] Searching for project: %s\n", projectName)
			searchResults, err := h.searchService.ExecuteSearchQuery(projectName)

			if err != nil {
				fmt.Printf("[ChatHandler] ERROR: Search failed: %v\n", err)
				chatResp.Content = "I couldn't find a project with that name. Please check the project name and try again."
			} else if len(searchResults) == 0 {
				fmt.Println("[ChatHandler] No matching projects found")
				chatResp.Content = "I couldn't find a project with that name. Please check the project name and try again."
			} else {
				fmt.Printf("[ChatHandler] Found %d matching projects\n", len(searchResults))

				// Get project ID from first result
				projectID, ok := searchResults[0]["projectId"].(string)
				if !ok || projectID == "" {
					fmt.Println("[ChatHandler] ERROR: Project ID missing in search result")
					chatResp.Content = "I found the project but couldn't generate a summary because of a data issue."
				} else {
					// Generate summary
					fmt.Printf("[ChatHandler] Generating summary for project %s\n", projectID)
					summary, err := h.summaryService.GenerateProjectSummary(projectID)

					if err != nil {
						fmt.Printf("[ChatHandler] ERROR: Summary generation failed: %v\n", err)
						chatResp.Content = "I found the project but couldn't generate a summary. Please try again later."
					} else {
						fmt.Println("[ChatHandler] Summary generated successfully")

						// Extract summary content
						if summaryContent, ok := summary["summary"].(string); ok {
							chatResp.Content = summaryContent
						} else {
							chatResp.Content = fmt.Sprintf("Here's what I found about the project: %s", searchResults[0]["projectName"].(string))
						}
					}
				}
			}
		} else {
			fmt.Println("[ChatHandler] No project name could be extracted")
			chatResp.Content = "Please specify which project you'd like a summary for."
		}
	} else if strings.Contains(messageContent, "search") ||
		strings.Contains(messageContent, "find") ||
		strings.Contains(messageContent, "projects about") ||
		strings.Contains(messageContent, "projects on") {
		fmt.Println("[ChatHandler] INTENT DETECTED: Search query")

		// Handle search intent - Use topic search for better semantic understanding
		query := extractSearchQuery(messageContent)
		fmt.Printf("[ChatHandler] Extracted search query: %s\n", query)

		fmt.Println("[ChatHandler] Executing topic search query")
		results, err := h.searchService.ExecuteTopicSearchQuery(query)

		if err != nil {
			fmt.Printf("[ChatHandler] ERROR: Topic search failed: %v\n", err)
			chatResp.Content = "I couldn't find any projects matching your query. Try a different search term."
		} else if len(results) == 0 {
			fmt.Println("[ChatHandler] No search results found")
			chatResp.Content = "I couldn't find any projects matching your query. Try a different search term."
		} else {
			fmt.Printf("[ChatHandler] Found %d search results\n", len(results))

			// Format search results for chat
			var responseText strings.Builder
			responseText.WriteString(fmt.Sprintf("I found %d projects matching your query:\n\n", len(results)))

			resultLimit := 5
			if len(results) < resultLimit {
				resultLimit = len(results)
			}

			for i := 0; i < resultLimit; i++ {
				projectName, _ := results[i]["projectName"].(string)
				projectDesc, _ := results[i]["description"].(string)
				if len(projectDesc) > 100 {
					projectDesc = projectDesc[:100] + "..."
				}

				responseText.WriteString(fmt.Sprintf("**%s**\n%s\n\n", projectName, projectDesc))
			}

			if len(results) > resultLimit {
				responseText.WriteString("... and more results.")
			}

			chatResp.Content = responseText.String()
		}
	} else {
		fmt.Println("[ChatHandler] INTENT DETECTED: General question - Using LLM")

		// For general questions, use the LLM directly
		fmt.Printf("[ChatHandler] Sending to LLM provider: %s\n", message)
		llmResponse, err := h.provider.Chat([]interfaces.Message{
			{Role: "user", Content: message},
		})

		if err != nil {
			fmt.Printf("[ChatHandler] ERROR: LLM request failed: %v\n", err)
			chatResp.Content = "I'm sorry, I couldn't process your request right now."
		} else {
			fmt.Println("[ChatHandler] LLM response received successfully")
			chatResp.Content = llmResponse.Content
		}
	}

	fmt.Println("[ChatHandler] Sending response to client")
	response.RespondWithJSON(w, http.StatusOK, chatResp)
}

// Extract search query from message
func extractSearchQuery(message string) string {
	message = strings.ToLower(message)
	message = strings.TrimPrefix(message, "search ")
	message = strings.TrimPrefix(message, "find ")
	message = strings.TrimPrefix(message, "projects about ")
	message = strings.TrimPrefix(message, "projects on ")
	return message
}

// Detect if the message is requesting a project summary
func matchesSummaryIntent(message string) bool {
	summaryTerms := []string{"summary", "summarize", "tell me about", "describe", "explain", "summarise"}
	projectTerms := []string{"project", "app", "application", "system", "platform"}

	hasSummaryTerm := false
	for _, term := range summaryTerms {
		if strings.Contains(message, term) {
			hasSummaryTerm = true
			break
		}
	}

	hasProjectTerm := false
	for _, term := range projectTerms {
		if strings.Contains(message, term) {
			hasProjectTerm = true
			break
		}
	}

	return hasSummaryTerm && hasProjectTerm
}

// Extract potential project name from a message
func extractProjectName(message string) string {
	// First try common patterns
	patterns := []string{
		`(?i)(?:summary|summarize|about|describe|explain).*?(?:project|app|system|platform)\s+(?:called|named)?\s*["']?([^"'.?!]+)["']?`,
		`(?i)(?:summary|summarize|about|describe|explain)\s+(?:the|a)?\s*["']?([^"'.?!]+)["']?(?:\s+project|\s+app|\s+system|\s+platform)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(message)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	// If no match with patterns, try extracting with keyword removal
	words := strings.Fields(message)
	filteredWords := make([]string, 0)

	// Remove common question words and filler words
	stopWords := map[string]bool{
		"summary": true, "summarize": true, "tell": true, "me": true, "about": true,
		"describe": true, "explain": true, "project": true, "app": true, "application": true,
		"the": true, "a": true, "an": true, "for": true, "of": true, "to": true,
		"give": true, "provide": true, "i": true, "want": true, "need": true,
		"would": true, "like": true, "get": true, "called": true, "named": true,
	}

	for _, word := range words {
		w := strings.ToLower(word)
		if !stopWords[w] {
			filteredWords = append(filteredWords, word)
		}
	}

	// Use the remaining words as potential project name
	if len(filteredWords) > 0 {
		return strings.Join(filteredWords, " ")
	}

	return ""
}

// GetChatHistory retrieves the chat history for a user
func (h *ChatHandler) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	// This would typically retrieve chat history from a database
	// For now, return an empty array
	response.RespondWithJSON(w, http.StatusOK, []ChatResponse{})
}
