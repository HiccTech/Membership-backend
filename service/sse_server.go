package service

import (
	"fmt"
	"hiccpet/service/middleware"
	"hiccpet/service/response"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	sseServerInstance *SSEServer
	once              sync.Once
)

type Client struct {
	id      string
	channel chan string
}

type SSEServer struct {
	clients map[string]*Client
	lock    sync.RWMutex
}

func NewSSEServer() *SSEServer {
	once.Do(func() {
		sseServerInstance = &SSEServer{
			clients: make(map[string]*Client),
		}
	})
	return sseServerInstance
}

func (s *SSEServer) AddClient(customerID string) *Client {
	s.lock.Lock()
	defer s.lock.Unlock()

	client := &Client{
		id:      customerID,
		channel: make(chan string, 10),
	}
	s.clients[customerID] = client
	return client
}

func (s *SSEServer) RemoveClient(customerID string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if client, ok := s.clients[customerID]; ok {
		close(client.channel)
		delete(s.clients, customerID)
	}
}

func (s *SSEServer) PushToClient(customerID, msg string) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if client, ok := s.clients[customerID]; ok {
		select {
		case client.channel <- msg:
		default:
			// 如果阻塞就丢弃
		}
	}
}

// Gin SSE Handler
func (s *SSEServer) Handler(c *gin.Context) {

	claims, err := middleware.VerifyShopifyToken(c.Query("token"))
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid token: "+err.Error())
		c.Abort()
		return
	}
	customerID := claims.Sub
	// customerID := c.Query("token")
	if customerID == "" {
		c.String(http.StatusBadRequest, "missing customer_id")
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "streaming unsupported")
		return
	}

	client := s.AddClient(customerID)
	defer s.RemoveClient(customerID)

	// 客户一连上就推送一条「已连接」的消息
	client.channel <- `{"code":0,"message":"connected"}`

	ctx := c.Request.Context()

	for {
		select {
		case msg := <-client.channel:
			fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
			flusher.Flush()
		case <-ctx.Done():
			return
		}
	}
}
