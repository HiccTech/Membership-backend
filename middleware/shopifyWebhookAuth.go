package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hiccpet/service/config"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 返回 HMAC-SHA256 的原始 bytes（未编码）
func computeShopifyHMAC(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func ShopifyWebhookAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取原始 body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
			return
		}

		// 读取后把 body 重置回去，供后续 handler 使用
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		c.Request.ContentLength = int64(len(body)) // 可选，保持 Content-Length 一致

		// 取得 header 并解 base64
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

		// 计算本地 hmac（原始 bytes）并比较
		expected := computeShopifyHMAC(body, []byte(config.Cfg.WebhookSecret))
		if !hmac.Equal(sig, expected) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}

		// 通过
		c.Next()
	}
}
