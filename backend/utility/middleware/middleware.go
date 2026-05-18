package middleware

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}

func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		msg := "OK"
		if len(c.Errors) > 0 {
			msg = c.Errors.Last().Error()
			c.JSON(http.StatusInternalServerError, Response{
				Message: msg,
			})
			return
		}

		res, exists := c.Get("response")
		if !exists {
			res = nil
		}
		c.JSON(http.StatusOK, Response{
			Message: msg,
			Data:    res,
		})
	}
}

type Response struct {
	Message string      `json:"message" dc:"消息提示"`
	Data    interface{} `json:"data"    dc:"执行结果"`
}
