package indexer

import (
	embedder2 "SuperBizAgent/internal/ai/embedder"
	"SuperBizAgent/utility/client"
	"SuperBizAgent/utility/common"
	"context"
	"encoding/json"

	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino/schema"
)

func NewMilvusIndexer(ctx context.Context) (*milvus.Indexer, error) {
	cli, err := client.NewMilvusClient(ctx)
	if err != nil {
		return nil, err
	}
	eb, err := embedder2.DoubaoEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config := &milvus.IndexerConfig{
		Client:     cli,
		Collection: common.MilvusCollectionName,
		Fields:     client.Fields,
		Embedding:  eb,
		// 默认 DocumentConverter 把向量转成 BinaryVector（[]byte），
		// 但 schema 是 FloatVector，需用 float32
		DocumentConverter: floatDocumentConverter,
	}
	indexer, err := milvus.NewIndexer(ctx, config)
	if err != nil {
		return nil, err
	}
	return indexer, nil
}

type indexRow struct {
	ID       string    `json:"id" milvus:"name:id"`
	Content  string    `json:"content" milvus:"name:content"`
	Vector   []float32 `json:"vector" milvus:"name:vector"`
	Metadata []byte    `json:"metadata" milvus:"name:metadata"`
}

func floatDocumentConverter(ctx context.Context, docs []*schema.Document, vectors [][]float64) ([]interface{}, error) {
	rows := make([]interface{}, len(docs))
	for i, doc := range docs {
		meta, err := json.Marshal(doc.MetaData)
		if err != nil {
			return nil, err
		}
		vec := make([]float32, len(vectors[i]))
		for j, v := range vectors[i] {
			vec[j] = float32(v)
		}
		rows[i] = &indexRow{
			ID:       doc.ID,
			Content:  doc.Content,
			Vector:   vec,
			Metadata: meta,
		}
	}
	return rows, nil
}
