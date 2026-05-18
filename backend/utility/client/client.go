package client

import (
	"SuperBizAgent/utility/common"
	"SuperBizAgent/utility/config"
	"context"
	"fmt"

	cli "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func NewMilvusClient(ctx context.Context) (cli.Client, error) {
	milvusAddr := config.GetString("milvus_addr")
	if milvusAddr == "" {
		milvusAddr = "localhost:19530"
	}
	// 1. 先连接default数据库
	defaultClient, err := cli.NewClient(ctx, cli.Config{
		Address: milvusAddr,
		DBName:  "default",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to default database: %w", err)
	}
	// 2. 检查agent数据库是否存在，不存在则创建
	databases, err := defaultClient.ListDatabases(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}
	agentDBExists := false
	for _, db := range databases {
		if db.Name == common.MilvusDBName {
			agentDBExists = true
			break
		}
	}
	if !agentDBExists {
		err = defaultClient.CreateDatabase(ctx, common.MilvusDBName)
		if err != nil {
			return nil, fmt.Errorf("failed to create agent database: %w", err)
		}
	}

	// 3. 创建连接到agent数据库的客户端
	agentClient, err := cli.NewClient(ctx, cli.Config{
		Address: milvusAddr,
		DBName:  common.MilvusDBName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to agent database: %w", err)
	}
	// 4. 检查biz collection是否存在，不存在则创建
	collections, err := agentClient.ListCollections(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}

	bizCollectionExists := false
	for _, collection := range collections {
		if collection.Name == common.MilvusCollectionName {
			bizCollectionExists = true
			break
		}
	}

	if !bizCollectionExists {
		// 创建biz collection的schema
		schema := &entity.Schema{
			CollectionName:    common.MilvusCollectionName,
			Description:       "Business knowledge collection",
			Fields:            Fields,
			EnableDynamicField: false,
		}

		err = agentClient.CreateCollection(ctx, schema, entity.DefaultShardNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to create biz collection: %w", err)
		}

		// 为id字段创建autoindex索引
		idIndex, err := entity.NewIndexAUTOINDEX(entity.L2)
		if err != nil {
			return nil, fmt.Errorf("failed to create id index: %w", err)
		}
		err = agentClient.CreateIndex(ctx, common.MilvusCollectionName, "id", idIndex, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create id index: %w", err)
		}

		// 为content字段创建autoindex索引
		contentIndex, err := entity.NewIndexAUTOINDEX(entity.L2)
		if err != nil {
			return nil, fmt.Errorf("failed to create content index: %w", err)
		}
		err = agentClient.CreateIndex(ctx, common.MilvusCollectionName, "content", contentIndex, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create content index: %w", err)
		}

		// 为vector字段创建autoindex索引
		vectorIndex, err := entity.NewIndexAUTOINDEX(entity.L2)
		if err != nil {
			return nil, fmt.Errorf("failed to create vector index: %w", err)
		}
		err = agentClient.CreateIndex(ctx, common.MilvusCollectionName, "vector", vectorIndex, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create vector index: %w", err)
		}
	}

	// 确保集合已加载到内存
	err = agentClient.LoadCollection(ctx, common.MilvusCollectionName, false)
	if err != nil {
		return nil, fmt.Errorf("failed to load collection: %w", err)
	}

	// 关闭default数据库连接
	defaultClient.Close()

	return agentClient, nil
}

// Fields 统一定义 Milvus biz collection 的 schema，供 indexer 和 retriever 共用
var Fields = []*entity.Field{
	{
		Name:     "id",
		DataType: entity.FieldTypeVarChar,
		TypeParams: map[string]string{
			"max_length": "256",
		},
		PrimaryKey: true,
	},
	{
		Name:     "vector",
		DataType: entity.FieldTypeFloatVector,
		TypeParams: map[string]string{
			"dim": "2048",
		},
	},
	{
		Name:     "content",
		DataType: entity.FieldTypeVarChar,
		TypeParams: map[string]string{
			"max_length": "8192",
		},
	},
	{
		Name:     "metadata",
		DataType: entity.FieldTypeJSON,
	},
}
