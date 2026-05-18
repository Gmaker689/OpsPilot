# OpsPilot - 智能 OnCall 助手

基于 AI Agent 的智能运维助手，提供日志诊断、知识库问答、指标查询等运维能力。

## 架构概览

```
OpsPilot/
├── backend/                  # Go 后端服务
│   ├── api/chat/            # HTTP 接口（对话、上传、AI Ops）
│   ├── internal/
│   │   ├── ai/
│   │   │   ├── agent/       # Agent 编排（ChatPipeline、Plan-Execute-Replan）
│   │   │   ├── embedder/    # 文本向量化（Doubao text-embedding-v4）
│   │   │   ├── indexer/     # Milvus 向量索引
│   │   │   ├── loader/      # 文档加载（PDF/TXT/MD/CSV/DOC/DOCX）
│   │   │   ├── retriever/   # Milvus 向量检索
│   │   │   └── tools/       # Agent 工具（SLS 日志、指标告警、MySQL、内部文档）
│   │   ├── controller/      # 路由控制器
│   │   └── logic/           # 业务逻辑（Chat、SSE）
│   ├── mcp-server/          # MCP 服务（阿里云 SLS 日志查询）
│   ├── utility/             # 工具库（配置、中间件、通用函数）
│   └── manifest/            # 配置与部署（Docker、config.yaml）
├── frontend/                # Web 前端
│   ├── index.html           # 智能对话界面
│   ├── app.js               # 前端逻辑（普通/流式对话、文件上传、历史记录）
│   └── styles.css           # 样式
└── docs/                    # 知识库文档
```

## 核心能力

- **智能对话**：普通请求-响应 / SSE 流式响应两种模式
- **RAG 知识库问答**：文档上传 → 分块 → 向量化 → Milvus 索引 → 检索增强生成
- **AI Ops 诊断**：Plan-Execute-Replan 架构，自主规划并执行多步运维任务
- **MCP 工具链**：SLS 日志查询、Prometheus 指标告警、MySQL 数据库操作、内部文档检索
- **ReAct Agent**：思考-行动-观察循环，自主调用工具完成诊断

## 技术栈


| 层级       | 技术                                                                       |
| ---------- | -------------------------------------------------------------------------- |
| 后端框架   | Go + Gin                                                                   |
| AI 编排    | [CloudWeGo Eino](https://github.com/cloudwego/eino) (Agent/Graph/Pipeline) |
| LLM        | Qwen (via SiliconFlow)                                                     |
| Embedding  | Qwen text-embedding-v4                                                     |
| 向量数据库 | Milvus                                                                     |
| 协议       | MCP (Model Context Protocol)                                               |
| 前端       | Vanilla JS + SSE + Marked                                                  |

## 快速开始

### 方式一：Docker Compose 一键部署

```bash
# 1. 配置环境变量
cp .env.example .env
# 编辑 .env 填写 LLM API Key 和阿里云 SLS 凭证

# 2. 启动全部服务
docker-compose up -d

# 3. 查看运行状态
docker-compose ps
```

服务启动后访问：


| 服务 | 地址 | 说明 |
| --- | --- | --- |
| 前端 | `http://localhost:8098` | Web 对话界面 |
| 后端 API | `http://localhost:8099` | REST API |
| Attu | `http://localhost:8000` | Milvus 管理界面 `--profile debug` |

可选启动 Attu：

```bash
docker-compose --profile debug up -d attu
```

### 方式二：手动启动（开发）

**环境依赖**：Go 1.26+、Milvus 实例

```bash
# 1. 启动 Milvus
cd backend/manifest/docker
docker-compose up -d etcd minio standalone

# 2. 启动 MCP Server（可选）
cd backend/mcp-server
cp .env.example .env
# 编辑 .env 填写阿里云凭证
go run main.go

# 3. 启动后端
cd backend
# 编辑 manifest/config/config.yaml 填写 API Key
go run main.go

# 4. 启动前端
cd frontend
./start.sh
```

### 使用

1. 打开前端页面，进入「智能 OnCall 助手」
2. 首次使用，通过下方工具栏上传知识库文档（支持 TXT/MD/Markdown）
3. 选择对话模式（快速 / 流式），开始对话
4. 点击「AI Ops」按钮进入运维诊断模式

## API 接口


| 方法 | 路径               | 说明             |
| ---- | ------------------ | ---------------- |
| POST | `/api/chat`        | 普通对话         |
| POST | `/api/chat_stream` | 流式对话 (SSE)   |
| POST | `/api/upload`      | 上传文件到知识库 |
| POST | `/api/ai_ops`      | AI Ops 诊断      |

## 配置说明

编辑 [backend/manifest/config/config.yaml](backend/manifest/config/config.yaml)：

- `ds_think_chat_model` — 深度思考模型（Plan 阶段）
- `ds_quick_chat_model` — 快速模型（Execute 阶段）
- `doubao_embedding_model` — 文本向量化模型
- `milvus_addr` — Milvus 服务地址
- `mcp_url` — MCP Server 地址
- `file_dir` — 知识库文件存储目录

## 文档

- [告警处理手册](docs/告警处理手册.md)
- [MCP Server 说明](backend/mcp-server/README.md)
- [前端使用指南](frontend/README.md)
