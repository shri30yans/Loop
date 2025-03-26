package ollama

import (
    "Loop_backend/internal/ai/interfaces"
)

// Ensure Provider implements the interfaces.Provider interface
var _ interfaces.Provider = (*Provider)(nil)
