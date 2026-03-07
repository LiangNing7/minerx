package text

import (
	"context"
	"encoding/json"

	"github.com/tmc/langchaingo/llms/ollama"

	"github.com/LiangNing7/goutils/pkg/log"

	"github.com/LiangNing7/minerx/internal/pkg/embedding/embedder/minerx"
)

// EmbeddingData holds the embedding data.
type EmbeddingData struct {
	// Common embedding data structure.
	Data minerx.EmbeddingData

	inputText string
}

// embedder implements the embedding interface for text data.
type embedder struct {
	// Client to interact with the embedding service.
	client *ollama.LLM
}

// NewEmbedder initializes a new embedder with the specified options.
func NewEmbedder(base string, client *ollama.LLM) *embedder {
	return &embedder{client: client}
}

// Embedding performs the embedding operation on the provided input.
func (emb *embedder) Embedding(ctx context.Context, input any) string {
	data, ok := input.(EmbeddingData)
	if !ok {
		log.W(ctx).Warnw("Invalid input type for embedding")
		return ""
	}

	if data.inputText == "" {
		log.W(ctx).Warnw("Encountered empty inputText when send embedding model request")
		return ""
	}

	embs, err := emb.client.CreateEmbedding(ctx, []string{data.inputText})
	if err != nil {
		log.W(ctx).Warnw("Failed to embed input text", "err", err)
		return ""
	}
	if len(embs) != 1 {
		log.W(ctx).Warnw("Embedding output is not equal to 1")
		return ""
	}

	// Store the embedding results.
	data.Data.Emb = embs[0]
	ret, _ := json.Marshal(data)
	return string(ret)
}
