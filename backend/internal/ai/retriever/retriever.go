package retriever

import (
	"SuperBizAgent/internal/ai/embedder"
	"SuperBizAgent/utility/client"
	"SuperBizAgent/utility/common"
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func NewMilvusRetriever(ctx context.Context) (rtr retriever.Retriever, err error) {
	cli, err := client.NewMilvusClient(ctx)
	if err != nil {
		return nil, err
	}
	eb, err := embedder.DoubaoEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	sp, _ := entity.NewIndexAUTOINDEXSearchParam(1)
	r, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:      cli,
		Collection:  common.MilvusCollectionName,
		VectorField: "vector",
		OutputFields: []string{"content", "metadata"},
		TopK:      1,
		Embedding:   eb,
		MetricType:  entity.L2,
		Sp:          sp,
		// 覆盖默认 BinaryVector 转换，Embedding 输出 float 向量需用 FloatVector
		VectorConverter: floatVectorConverter,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func floatVectorConverter(ctx context.Context, vectors [][]float64) ([]entity.Vector, error) {
	if len(vectors) != 1 {
		return nil, fmt.Errorf("expected 1 vector, got %d", len(vectors))
	}
	vec := make([]float32, len(vectors[0]))
	for i, v := range vectors[0] {
		vec[i] = float32(v)
	}
	return []entity.Vector{entity.FloatVector(vec)}, nil
}
