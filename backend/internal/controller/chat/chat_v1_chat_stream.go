package chat

import (
	v1 "SuperBizAgent/api/chat/v1"
	"SuperBizAgent/internal/ai/agent/chat_pipeline"
	"SuperBizAgent/utility/log_call_back"
	"SuperBizAgent/utility/mem"
	"io"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

func (c *ControllerV1) ChatStream(ctx *gin.Context) {
	var req v1.ChatStreamReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	client, err := c.service.Create(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	userMessage := &chat_pipeline.UserMessage{
		ID:      req.Id,
		Query:   req.Question,
		History: mem.GetSimpleMemory(req.Id).GetMessages(),
	}

	runner, err := chat_pipeline.BuildChatAgent(ctx.Request.Context())
	if err != nil {
		client.SendToClient("error", err.Error())
		return
	}

	sr, err := runner.Stream(ctx.Request.Context(), userMessage, compose.WithCallbacks(log_call_back.LogCallback(nil)))
	if err != nil {
		client.SendToClient("error", err.Error())
		return
	}
	defer sr.Close()

	var fullResponse strings.Builder

	defer func() {
		completeResponse := fullResponse.String()
		if completeResponse != "" {
			mem.GetSimpleMemory(req.Id).SetMessages(schema.UserMessage(req.Question))
			mem.GetSimpleMemory(req.Id).SetMessages(schema.SystemMessage(completeResponse))
		}
	}()

	for {
		chunk, err := sr.Recv()
		if err == io.EOF {
			client.SendToClient("done", "Stream completed")
			return
		}
		if err != nil {
			client.SendToClient("error", err.Error())
			return
		}
		fullResponse.WriteString(chunk.Content)
		client.SendToClient("message", chunk.Content)
	}
}
