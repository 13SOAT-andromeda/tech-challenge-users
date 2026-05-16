package handlers

import (
	"net/http"
	"testing"

	"tech-challenge-users/internal/application/services"

	"github.com/gin-gonic/gin"
)

func buildCompanyTestRouter() *gin.Engine {
	cr := newMockCompanyRepo()
	svc := services.NewCompanyService(cr)
	h := NewCompanyHandler(svc)

	return setupRouter(func(g *gin.RouterGroup) {
		g.POST("/companies", h.Create)
		g.GET("/companies/:id", h.GetByID)
		g.PUT("/companies/:id", h.Update)
		g.DELETE("/companies/:id", h.Delete)
	})
}

const validCompanyBody = `{"name":"ACME","email":"acme@acme.com","document":"11222333000181","contact":"11999999999"}`

func TestCompanyHandler_Create_Success(t *testing.T) {
	r := buildCompanyTestRouter()
	w := doRequest(r, http.MethodPost, "/companies", validCompanyBody)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCompanyHandler_Create_InvalidJSON(t *testing.T) {
	r := buildCompanyTestRouter()
	w := doRequest(r, http.MethodPost, "/companies", `{bad}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCompanyHandler_Create_InvalidDocument(t *testing.T) {
	r := buildCompanyTestRouter()
	body := `{"name":"X","email":"x@x.com","document":"00000000000000","contact":"11"}`
	w := doRequest(r, http.MethodPost, "/companies", body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCompanyHandler_Create_DuplicateDocument(t *testing.T) {
	r := buildCompanyTestRouter()
	doRequest(r, http.MethodPost, "/companies", validCompanyBody)
	w := doRequest(r, http.MethodPost, "/companies", validCompanyBody)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestCompanyHandler_Create_WithAddress(t *testing.T) {
	r := buildCompanyTestRouter()
	body := `{"name":"ACME","email":"acme@acme.com","document":"11222333000181","contact":"11","address":{"street":"Av Paulista","number":"1","city":"SP","country":"BR","zip_code":"01310100"}}`
	w := doRequest(r, http.MethodPost, "/companies", body)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestCompanyHandler_GetByID_Success(t *testing.T) {
	r := buildCompanyTestRouter()
	doRequest(r, http.MethodPost, "/companies", validCompanyBody)
	w := doRequest(r, http.MethodGet, "/companies/1", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCompanyHandler_GetByID_NotFound(t *testing.T) {
	r := buildCompanyTestRouter()
	w := doRequest(r, http.MethodGet, "/companies/999", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCompanyHandler_GetByID_InvalidID(t *testing.T) {
	r := buildCompanyTestRouter()
	w := doRequest(r, http.MethodGet, "/companies/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCompanyHandler_Update_Success(t *testing.T) {
	r := buildCompanyTestRouter()
	doRequest(r, http.MethodPost, "/companies", validCompanyBody)
	body := `{"name":"ACME Updated","contact":"11988888888"}`
	w := doRequest(r, http.MethodPut, "/companies/1", body)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCompanyHandler_Update_WithAddress(t *testing.T) {
	r := buildCompanyTestRouter()
	doRequest(r, http.MethodPost, "/companies", validCompanyBody)
	body := `{"address":{"street":"Rua Nova","number":"99","city":"BH"}}`
	w := doRequest(r, http.MethodPut, "/companies/1", body)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCompanyHandler_Update_NotFound(t *testing.T) {
	r := buildCompanyTestRouter()
	w := doRequest(r, http.MethodPut, "/companies/999", `{"name":"X"}`)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCompanyHandler_Update_InvalidID(t *testing.T) {
	r := buildCompanyTestRouter()
	w := doRequest(r, http.MethodPut, "/companies/abc", `{"name":"X"}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCompanyHandler_Delete_Success(t *testing.T) {
	r := buildCompanyTestRouter()
	doRequest(r, http.MethodPost, "/companies", validCompanyBody)
	w := doRequest(r, http.MethodDelete, "/companies/1", "")
	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestCompanyHandler_Delete_NotFound(t *testing.T) {
	r := buildCompanyTestRouter()
	w := doRequest(r, http.MethodDelete, "/companies/999", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCompanyHandler_Delete_InvalidID(t *testing.T) {
	r := buildCompanyTestRouter()
	w := doRequest(r, http.MethodDelete, "/companies/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
