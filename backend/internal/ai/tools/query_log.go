package tools

import (
	"SuperBizAgent/utility/config"
	"context"

	e_mcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetLogMcpTool() ([]tool.BaseTool, error) {
	mcpUrl := config.GetString("mcp_url")
	ctx := context.Background()
	cli, err := client.NewStreamableHttpClient(mcpUrl)
	if err != nil {
		return []tool.BaseTool{}, err
	}
	err = cli.Start(ctx)
	if err != nil {
		return []tool.BaseTool{}, err
	}
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "example-client",
		Version: "1.0.0",
	}
	if _, err = cli.Initialize(ctx, initRequest); err != nil {
		return []tool.BaseTool{}, err
	}
	mcpTools, err := e_mcp.GetTools(ctx, &e_mcp.Config{Cli: cli})
	if err != nil {
		return []tool.BaseTool{}, err
	}
	return mcpTools, nil
}
