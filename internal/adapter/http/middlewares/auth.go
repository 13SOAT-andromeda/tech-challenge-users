package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID    = "userID"
	ContextUserEmail = "userEmail"
	ContextUserRole  = "userRole"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
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
