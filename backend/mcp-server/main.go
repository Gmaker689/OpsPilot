package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type config struct {
	AccessKeyID     string
	AccessKeySecret string
	Endpoint        string
	Project         string
	Logstore        string
	ListenAddr      string
}

func loadConfig() *config {
	return &config{
		AccessKeyID:     os.Getenv("ALIBABA_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("ALIBABA_ACCESS_KEY_SECRET"),
		Endpoint:        os.Getenv("SLS_ENDPOINT"),
		Project:         os.Getenv("SLS_PROJECT"),
		Logstore:        os.Getenv("SLS_LOGSTORE"),
		ListenAddr:      getEnv("LISTEN_ADDR", ":8080"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := loadConfig()

	slsClient := &sls.Client{
		Endpoint:        cfg.Endpoint,
		AccessKeyID:     cfg.AccessKeyID,
		AccessKeySecret: cfg.AccessKeySecret,
	}

	mcpServer := server.NewMCPServer("ops-pilot-sls", "1.0.0")

	mcpServer.AddTool(mcp.Tool{
		Name:        "search_sls_logs",
		Description: "Search logs from Alibaba Cloud Simple Log Service (SLS). Use this tool to query log entries by keyword within a specified time range. This is the primary tool for diagnosing service issues, finding error logs, and investigating incidents.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"keyword": map[string]any{
					"type":        "string",
					"description": "Search keyword for full-text search in logs. Supports SLS query syntax (e.g., 'panic', 'error', 'response', 'region AND mismatch').",
				},
				"start_time": map[string]any{
					"type":        "string",
					"description": "Start time for log search in RFC3339 format (e.g., '2026-05-17T09:00:00+08:00'). Defaults to 1 hour before end_time.",
				},
				"end_time": map[string]any{
					"type":        "string",
					"description": "End time for log search in RFC3339 format (e.g., '2026-05-17T10:00:00+08:00'). Defaults to current time.",
				},
				"limit": map[string]any{
					"type":        "number",
					"description": "Maximum number of log entries to return. Default: 20, Max: 100.",
				},
			},
			Required: []string{"keyword"},
		},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		keyword := req.GetString("keyword", "")
		startTimeStr := req.GetString("start_time", "")
		endTimeStr := req.GetString("end_time", "")
		limit := req.GetInt("limit", 20)
		if limit > 100 {
			limit = 100
		}

		endTime := time.Now()
		if endTimeStr != "" {
			parsed, err := time.Parse(time.RFC3339, endTimeStr)
			if err != nil {
				return errorResult(fmt.Sprintf("invalid end_time format: %v, use RFC3339", err)), nil
			}
			endTime = parsed
		}

		startTime := endTime.Add(-1 * time.Hour)
		if startTimeStr != "" {
			parsed, err := time.Parse(time.RFC3339, startTimeStr)
			if err != nil {
				return errorResult(fmt.Sprintf("invalid start_time format: %v, use RFC3339", err)), nil
			}
			startTime = parsed
		}

		log.Printf("SLS search: project=%s logstore=%s keyword=%q from=%s to=%s limit=%d",
			cfg.Project, cfg.Logstore, keyword, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339), limit)

		searchResp, err := slsClient.GetLogsV2(cfg.Project, cfg.Logstore, &sls.GetLogRequest{
			From:    startTime.Unix(),
			To:      endTime.Unix(),
			Query:   keyword,
			Lines:   int64(limit),
			Reverse: true,
		})
		if err != nil {
			log.Printf("SLS search error: %v", err)
			return errorResult(fmt.Sprintf("SLS search failed: %v", err)), nil
		}

		logs := searchResp.Logs
		if logs == nil {
			logs = []map[string]string{}
		}

		result := map[string]any{
			"success": true,
			"count":   len(logs),
			"logs":    logs,
			"query": map[string]any{
				"keyword":    keyword,
				"start_time": startTime.Format(time.RFC3339),
				"end_time":   endTime.Format(time.RFC3339),
				"project":    cfg.Project,
				"logstore":   cfg.Logstore,
			},
		}

		jsonBytes, _ := json.MarshalIndent(result, "", "  ")
		log.Printf("SLS search completed: %d entries found", len(logs))
		return textResult(string(jsonBytes)), nil
	})

	mcpServer.AddTool(mcp.Tool{
		Name:        "get_current_time",
		Description: "Get the current server time. Use this to get an accurate RFC3339 timestamp for calculating time ranges when searching logs.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]any{},
		},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		now := time.Now()
		result := map[string]any{
			"current_time":  now.Format(time.RFC3339),
			"timestamp_unix": now.Unix(),
			"timezone":      now.Format("MST"),
			"one_hour_ago":  now.Add(-1 * time.Hour).Format(time.RFC3339),
		}
		jsonBytes, _ := json.MarshalIndent(result, "", "  ")
		return textResult(string(jsonBytes)), nil
	})

	httpServer := server.NewStreamableHTTPServer(mcpServer)

	log.Printf("MCP SLS Server starting on %s", cfg.ListenAddr)
	log.Printf("SLS Endpoint: %s, Project: %s, Logstore: %s", cfg.Endpoint, cfg.Project, cfg.Logstore)
	if err := httpServer.Start(cfg.ListenAddr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: text},
		},
	}
}

func errorResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: msg},
		},
		IsError: true,
	}
}
