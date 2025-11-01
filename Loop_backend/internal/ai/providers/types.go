package providers

import (
    "Loop_backend/internal/ai/interfaces"
)

// Ensure backward compatibility
type (
    Provider = interfaces.Provider
    Message = interfaces.Message
    ChatResponse = interfaces.ChatResponse
)
