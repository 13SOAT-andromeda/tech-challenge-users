package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID    = "userID"
	ContextUserEmail = "userEmail"
	ContextUserRole  = "userRole"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for internal S2S token first
		internalToken := c.GetHeader("X-Internal-Token")
		expectedToken := os.Getenv("INTERNAL_AUTH_TOKEN")
		if internalToken != "" && expectedToken != "" && internalToken == expectedToken {
			c.Set(ContextUserID, os.Getenv("ADMIN_DOCUMENT"))
			c.Set(ContextUserEmail, os.Getenv("ADMIN_EMAIL"))
			c.Set(ContextUserRole, "administrator")
			c.Next()
			return
		}

		userID := c.GetHeader("X-User-Id")
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set(ContextUserID, userID)
		c.Set(ContextUserEmail, c.GetHeader("X-User-Email"))
		c.Set(ContextUserRole, c.GetHeader("X-User-Role"))

		c.Next()
	}
}
