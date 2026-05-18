package sse

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Client struct {
	Id          string
	Writer      gin.ResponseWriter
	messageChan chan string
}

type Service struct {
	clients sync.Map
}

func New() *Service {
	return &Service{}
}

func (s *Service) Create(c *gin.Context) (*Client, error) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	clientId := c.Query("client_id")
	if clientId == "" {
		clientId = uuid.New().String()
	}
	client := &Client{
		Id:          clientId,
		Writer:      c.Writer,
		messageChan: make(chan string, 100),
	}

	c.Writer.WriteString(fmt.Sprintf("id: %s\n", clientId))
	c.Writer.WriteString("event: connected\n")
	c.Writer.WriteString(fmt.Sprintf("data: {\"status\": \"connected\", \"client_id\": \"%s\"}\n\n", clientId))
	c.Writer.Flush()

	return client, nil
}

func (c *Client) SendToClient(eventType, data string) bool {
	msg := fmt.Sprintf(
		"id: %d\nevent: %s\ndata: %s\n\n",
		time.Now().UnixNano(), eventType, data,
	)
	c.Writer.WriteString(msg)
	c.Writer.Flush()
	return true
}
