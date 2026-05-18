package main

import (
	"SuperBizAgent/internal/controller/chat"
	"SuperBizAgent/utility/common"
	"SuperBizAgent/utility/config"
	"SuperBizAgent/utility/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("config init failed: %v", err)
	}

	common.FileDir = config.GetString("file_dir")

	r := gin.Default()

	api := r.Group("/api")
	api.Use(middleware.CORSMiddleware())
	api.Use(middleware.ResponseMiddleware())

	ctrl := chat.NewV1()
	api.POST("/chat", ctrl.Chat)
	api.POST("/chat_stream", ctrl.ChatStream)
	api.POST("/upload", ctrl.FileUpload)
	api.POST("/ai_ops", ctrl.AIOps)

	r.Run(":8099")
}
