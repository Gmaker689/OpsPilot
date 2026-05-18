package chat

import "github.com/gin-gonic/gin"

type IChatV1 interface {
	Chat(c *gin.Context)
	ChatStream(c *gin.Context)
	FileUpload(c *gin.Context)
	AIOps(c *gin.Context)
}
