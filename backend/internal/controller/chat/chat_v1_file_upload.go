package chat

import (
	v1 "SuperBizAgent/api/chat/v1"
	"SuperBizAgent/internal/ai/agent/knowledge_index_pipeline"
	loader2 "SuperBizAgent/internal/ai/loader"
	"SuperBizAgent/utility/client"
	"SuperBizAgent/utility/common"
	"SuperBizAgent/utility/log_call_back"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
	"github.com/gin-gonic/gin"
)

func (c *ControllerV1) FileUpload(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.Error(fmt.Errorf("请上传文件: %w", err))
		return
	}
	defer file.Close()

	if err := os.MkdirAll(common.FileDir, 0755); err != nil {
		ctx.Error(fmt.Errorf("创建目录失败 %s: %w", common.FileDir, err))
		return
	}

	savePath := filepath.Join(common.FileDir, header.Filename)
	dst, err := os.Create(savePath)
	if err != nil {
		ctx.Error(fmt.Errorf("保存文件失败: %w", err))
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		ctx.Error(fmt.Errorf("保存文件失败: %w", err))
		return
	}

	fileInfo, err := os.Stat(savePath)
	if err != nil {
		ctx.Error(fmt.Errorf("获取文件信息失败: %w", err))
		return
	}

	res := &v1.FileUploadRes{
		FileName: header.Filename,
		FilePath: savePath,
		FileSize: fileInfo.Size(),
	}

	if err := buildIntoIndex(ctx.Request.Context(), savePath); err != nil {
		ctx.Error(fmt.Errorf("构建知识库失败: %w", err))
		return
	}

	ctx.Set("response", res)
}

func buildIntoIndex(ctx context.Context, path string) error {
	r, err := knowledge_index_pipeline.BuildKnowledgeIndexing(ctx)
	if err != nil {
		return err
	}
	loader, err := loader2.NewFileLoader(ctx)
	if err != nil {
		return err
	}
	docs, err := loader.Load(ctx, document.Source{URI: path})
	if err != nil {
		return err
	}
	cli, err := client.NewMilvusClient(ctx)
	if err != nil {
		return err
	}
	source, _ := docs[0].MetaData["_source"].(string)
	source = strings.ReplaceAll(source, `\`, `\\`)
	expr := fmt.Sprintf(`metadata["_source"] == "%s"`, source)
	queryResult, err := cli.Query(ctx, common.MilvusCollectionName, []string{}, expr, []string{"id"})
	if err != nil {
		return err
	} else if len(queryResult) > 0 {
		var idsToDelete []string
		for _, column := range queryResult {
			if column.Name() == "id" {
				for i := 0; i < column.Len(); i++ {
					id, err := column.GetAsString(i)
					if err == nil {
						idsToDelete = append(idsToDelete, id)
					}
				}
			}
		}
		if len(idsToDelete) > 0 {
			deleteExpr := fmt.Sprintf(`id in ["%s"]`, strings.Join(idsToDelete, `","`))
			err = cli.Delete(ctx, common.MilvusCollectionName, "", deleteExpr)
			if err != nil {
				fmt.Printf("[warn] delete existing data failed: %v\n", err)
			} else {
				fmt.Printf("[info] deleted %d existing records with _source: %s\n", len(idsToDelete), docs[0].MetaData["_source"])
			}
		}
	}
	ids, err := r.Invoke(ctx, document.Source{URI: path}, compose.WithCallbacks(log_call_back.LogCallback(nil)))
	if err != nil {
		return fmt.Errorf("invoke index graph failed: %w", err)
	}
	fmt.Printf("[done] indexing file: %s, len of parts: %d\n", path, len(ids))
	return nil
}
