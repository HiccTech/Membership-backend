package middleware

import (
	"errors"
	"fmt"
	"hiccpet/service/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ShopifyClaims 对应 Shopify JWT 的标准字段
type ShopifyClaims struct {
	Sub  string `json:"sub"`  // Shopify Customer ID
	Dest string `json:"dest"` // 店铺域名
	jwt.RegisteredClaims
}

var ShopifyAppSecret = []byte("6e066851be38870b7638fc76b7e523d0")

func ShopifySessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			response.Error(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			response.Error(c, http.StatusUnauthorized, "invalid authorization header")
			c.Abort()
			return
		}
		claims, err := VerifyShopifyToken(tokenStr)
		if err != nil {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			response.Error(c, http.StatusUnauthorized, "invalid token: "+err.Error())
			c.Abort()
			return
		}
		fmt.Printf("claims: %+v\n", claims)
		// token 验证成功，把 claims 存到 context
		c.Set("shopifyClaims", claims)
		c.Next()
	}
}

// 验证 Shopify session token
func VerifyShopifyToken(tokenStr string) (*ShopifyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &ShopifyClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return ShopifyAppSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ShopifyClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
