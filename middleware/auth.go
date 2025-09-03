package middleware

import (
	"hiccpet/service/response"
	"hiccpet/service/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			response.Error(c, http.StatusUnauthorized, "Missing token")
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := utils.ParseToken(tokenString)
		if err != nil || !token.Valid {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			response.Error(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		c.Next()
	}
}
