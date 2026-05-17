package handlers

import (
	"net/http"
	"testing"

	"tech-challenge-users/internal/application/services"
	"tech-challenge-users/internal/domain"

	"github.com/gin-gonic/gin"
)

func buildEmployeeTestRouter() *gin.Engine {
	er := newMockEmployeeRepo()
	svc := services.NewEmployeeService(er)
	h := NewEmployeeHandler(svc)

	return setupRouter(func(g *gin.RouterGroup) {
		g.GET("/employees/:id", h.GetByID)
	})
}

func buildSeededEmployeeTestRouter() (*gin.Engine, *mockEmployeeRepo) {
	er := newMockEmployeeRepo()
	_ = er.Create(&domain.Employee{Position: "mechanic", PersonID: 10})

	svc := services.NewEmployeeService(er)
	h := NewEmployeeHandler(svc)

	r := setupRouter(func(g *gin.RouterGroup) {
		g.GET("/employees/:id", h.GetByID)
	})
	return r, er
}

func TestEmployeeHandler_GetByID_Success(t *testing.T) {
	r, _ := buildSeededEmployeeTestRouter()
	w := doRequest(r, http.MethodGet, "/employees/1", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestEmployeeHandler_GetByID_NotFound(t *testing.T) {
	r := buildEmployeeTestRouter()
	w := doRequest(r, http.MethodGet, "/employees/999", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestEmployeeHandler_GetByID_InvalidID(t *testing.T) {
	r := buildEmployeeTestRouter()
	w := doRequest(r, http.MethodGet, "/employees/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
