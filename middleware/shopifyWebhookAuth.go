package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hiccpet/service/config"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// 简单内存存储已处理的 Event ID（生产环境推荐用 Redis 或数据库）
var processedEvents = struct {
	sync.RWMutex
	m map[string]bool
}{m: make(map[string]bool)}

// Shopify HMAC 验证
func computeShopifyHMAC(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func ShopifyWebhookAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
			return
		}

		// 读取后重置 body
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		// 校验 HMAC
		shopifySig := c.GetHeader("X-Shopify-Hmac-Sha256")
		if shopifySig == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing signature"})
			return
		}
		sig, err := base64.StdEncoding.DecodeString(shopifySig)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid signature encoding"})
			return
		}
		expected := computeShopifyHMAC(body, []byte(config.Cfg.WebhookSecret))
		if !hmac.Equal(sig, expected) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}

		// 幂等性检查
		eventID := c.GetHeader("X-Shopify-Event-Id")
		if eventID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing event ID"})
			return
		}

		processedEvents.RLock()
		_, exists := processedEvents.m[eventID]
		processedEvents.RUnlock()

		if exists {
			// 已处理过，直接返回 200 OK
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "duplicate event"})
			return
		}

		// 标记为已处理
		processedEvents.Lock()
		processedEvents.m[eventID] = true
		processedEvents.Unlock()

		c.Next()
	}
}
