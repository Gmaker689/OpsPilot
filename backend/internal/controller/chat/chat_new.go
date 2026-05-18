package chat

import (
	"SuperBizAgent/api/chat"
	"SuperBizAgent/internal/logic/sse"
)

type ControllerV1 struct {
	service *sse.Service
}

func NewV1() chat.IChatV1 {
	return &ControllerV1{
		service: sse.New(),
	}
}
