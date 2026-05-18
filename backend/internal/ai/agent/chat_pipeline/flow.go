package chat_pipeline

import (
	"SuperBizAgent/internal/ai/tools"
	"context"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

func newReactAgentLambda(ctx context.Context) (lba *compose.Lambda, err error) {
	config := &react.AgentConfig{
		MaxStep:            25,
		ToolReturnDirectly: map[string]struct{}{}}
	chatModelIns11, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	config.ToolCallingModel = chatModelIns11

	var toolList []tool.BaseTool

	mcpTool, err := tools.GetLogMcpTool()
	if err != nil {
		log.Printf("warn: MCP tools unavailable: %v", err)
	} else {
		toolList = append(toolList, mcpTool...)
	}

	alertsTool, err := tools.NewPrometheusAlertsQueryTool()
	if err != nil {
		return nil, err
	}
	toolList = append(toolList, alertsTool)

	mysqlTool, err := tools.NewMysqlCrudTool()
	if err != nil {
		return nil, err
	}
	toolList = append(toolList, mysqlTool)

	timeTool, err := tools.NewGetCurrentTimeTool()
	if err != nil {
		return nil, err
	}
	toolList = append(toolList, timeTool)

	docsTool, err := tools.NewQueryInternalDocsTool()
	if err != nil {
		return nil, err
	}
	toolList = append(toolList, docsTool)

	config.ToolsConfig.Tools = toolList

	ins, err := react.NewAgent(ctx, config)
	if err != nil {
		return nil, err
	}
	lba, err = compose.AnyLambda(ins.Generate, ins.Stream, nil, nil)
	if err != nil {
		return nil, err
	}
	return lba, nil
}
