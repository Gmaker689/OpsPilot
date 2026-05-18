package embedder

import (
	"SuperBizAgent/utility/config"
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/cloudwego/eino/components/embedding"
)

func DoubaoEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	model := config.GetString("doubao_embedding_model.model")
	apiKey := config.GetString("doubao_embedding_model.api_key")
	dim := 2048
	embedder, err := dashscope.NewEmbedder(ctx, &dashscope.EmbeddingConfig{
		Model:      model,
		APIKey:     apiKey,
		Dimensions: &dim,
	})
	if err != nil {
		log.Printf("new embedder error: %v\n", err)
		return nil, err
	}
	return embedder, nil
}
