# OpsPilot MCP Server

基于 MCP (Model Context Protocol) 的阿里云 SLS 日志查询服务，为 OpsPilot AI Agent 提供日志诊断能力。

## 工具

| 工具 | 说明 |
|------|------|
| `search_sls_logs` | 按关键字和时间范围查询阿里云 SLS 日志 |
| `get_current_time` | 获取服务端当前时间，用于时间范围计算 |

## 快速启动

```bash
# 1. 配置
cp .env.example .env
# 编辑 .env 填写阿里云 AccessKey 和 SLS 信息

# 2. 启动
go run main.go

# 3. 编译部署
go build -o mcp-server .
./mcp-server
```

## 环境变量

| 变量 | 说明 |
|------|------|
| `ALIBABA_ACCESS_KEY_ID` | 阿里云 AccessKey ID |
| `ALIBABA_ACCESS_KEY_SECRET` | 阿里云 AccessKey Secret |
| `SLS_ENDPOINT` | SLS 地域域名 |
| `SLS_PROJECT` | SLS 项目名 |
| `SLS_LOGSTORE` | SLS 日志库名 |
| `LISTEN_ADDR` | 监听地址，默认 `:8080` |

## 接入 OpsPilot

修改主项目 `config.yaml` 中的 `mcp_url` 指向此服务：

```yaml
mcp_url: "http://localhost:8080/mcp"
```
