package chat

import (
	v1 "SuperBizAgent/api/chat/v1"
	"SuperBizAgent/internal/ai/agent/chat_pipeline"
	"SuperBizAgent/utility/log_call_back"
	"SuperBizAgent/utility/mem"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

func (c *ControllerV1) Chat(ctx *gin.Context) {
	var req v1.ChatReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
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
		ctx.Error(err)
		return
	}

	out, err := runner.Invoke(ctx.Request.Context(), userMessage, compose.WithCallbacks(log_call_back.LogCallback(nil)))
	if err != nil {
		ctx.Error(err)
		return
	}

	mem.GetSimpleMemory(req.Id).SetMessages(schema.UserMessage(req.Question))
	mem.GetSimpleMemory(req.Id).SetMessages(schema.AssistantMessage(out.Content, nil))

	ctx.Set("response", &v1.ChatRes{Answer: out.Content})
}
