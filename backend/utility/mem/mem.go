package mem

import (
	"sync"
	"time"

	"github.com/cloudwego/eino/schema"
)

const (
	defaultTTL       = 30 * time.Minute
	cleanupInterval  = 5 * time.Minute
)

var SimpleMemoryMap = make(map[string]*SimpleMemory)
var mu sync.Mutex

func init() {
	go func() {
		for {
			time.Sleep(cleanupInterval)
			cleanupExpired()
		}
	}()
}

func cleanupExpired() {
	mu.Lock()
	defer mu.Unlock()
	now := time.Now()
	for id, mem := range SimpleMemoryMap {
		mem.mu.Lock()
		expired := now.Sub(mem.lastAccess) > defaultTTL
		mem.mu.Unlock()
		if expired {
			delete(SimpleMemoryMap, id)
		}
	}
}

func GetSimpleMemory(id string) *SimpleMemory {
	mu.Lock()
	defer mu.Unlock()
	if mem, ok := SimpleMemoryMap[id]; ok {
		mem.mu.Lock()
		mem.lastAccess = time.Now()
		mem.mu.Unlock()
		return mem
	}
	newMem := &SimpleMemory{
		ID:            id,
		Messages:      []*schema.Message{},
		MaxWindowSize: 6,
		lastAccess:    time.Now(),
	}
	SimpleMemoryMap[id] = newMem
	return newMem
}

type SimpleMemory struct {
	ID            string            `json:"id"`
	Messages      []*schema.Message `json:"messages"`
	MaxWindowSize int
	lastAccess    time.Time
	mu            sync.Mutex
}

func (c *SimpleMemory) SetMessages(msg *schema.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastAccess = time.Now()
	c.Messages = append(c.Messages, msg)
	if len(c.Messages) > c.MaxWindowSize {
		excess := len(c.Messages) - c.MaxWindowSize
		if excess%2 != 0 {
			excess++
		}
		c.Messages = c.Messages[excess:]
	}
}

func (c *SimpleMemory) GetMessages() []*schema.Message {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastAccess = time.Now()
	return c.Messages
}
