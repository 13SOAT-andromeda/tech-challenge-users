package handlers

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealth(t *testing.T) {
	r := setupRouter(func(g *gin.RouterGroup) {
		g.GET("/health", Health)
	})
	w := doRequest(r, http.MethodGet, "/health", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
