package plan_execute_replan

import (
	"SuperBizAgent/internal/ai/models"
	"SuperBizAgent/internal/ai/tools"
	"context"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/planexecute"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

func NewExecutor(ctx context.Context) (adk.Agent, error) {
	var toolList []tool.BaseTool
	// log
	mcpTool, err := tools.GetLogMcpTool()
	if err != nil {
		log.Printf("warn: MCP tools unavailable: %v", err)
	} else {
		toolList = append(toolList, mcpTool...)
	}
	// alerts
	alertsTool, err := tools.NewPrometheusAlertsQueryTool()
	if err != nil {
		return nil, err
	}
	toolList = append(toolList, alertsTool)
	// file
	docsTool, err := tools.NewQueryInternalDocsTool()
	if err != nil {
		return nil, err
	}
	toolList = append(toolList, docsTool)
	// time
	timeTool, err := tools.NewGetCurrentTimeTool()
	if err != nil {
		return nil, err
	}
	toolList = append(toolList, timeTool)
	execModel, err := models.OpenAIForDeepSeekV3Quick(ctx)
	if err != nil {
		return nil, err
	}
	return planexecute.NewExecutor(ctx, &planexecute.ExecutorConfig{
		Model: execModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: toolList,
			},
		},
		MaxIterations: 999999,
	})
}
