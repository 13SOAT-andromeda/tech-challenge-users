package handlers

import (
	"net/http"
	"testing"

	"tech-challenge-users/internal/application/services"

	"github.com/gin-gonic/gin"
)

func buildVehicleTestRouter() *gin.Engine {
	vr := newMockVehicleRepo()
	svc := services.NewVehicleService(vr)
	h := NewVehicleHandler(svc)

	return setupRouter(func(g *gin.RouterGroup) {
		g.POST("/vehicles", h.Create)
		g.GET("/vehicles", h.List)
		g.GET("/vehicles/:id", h.GetByID)
		g.PUT("/vehicles/:id", h.Update)
		g.DELETE("/vehicles/:id", h.Delete)
	})
}

const validVehicleBody = `{"plate":"ABC1234","name":"Civic","year":2020,"brand":"Honda","color":"Prata"}`

func TestVehicleHandler_Create_Success(t *testing.T) {
	r := buildVehicleTestRouter()
	w := doRequest(r, http.MethodPost, "/vehicles", validVehicleBody)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestVehicleHandler_Create_InvalidJSON(t *testing.T) {
	r := buildVehicleTestRouter()
	w := doRequest(r, http.MethodPost, "/vehicles", `{bad}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestVehicleHandler_Create_InvalidPlate(t *testing.T) {
	r := buildVehicleTestRouter()
	body := `{"plate":"INVALID","name":"Civic","year":2020,"brand":"Honda","color":"Prata"}`
	w := doRequest(r, http.MethodPost, "/vehicles", body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestVehicleHandler_Create_DuplicatePlate(t *testing.T) {
	r := buildVehicleTestRouter()
	doRequest(r, http.MethodPost, "/vehicles", validVehicleBody)
	w := doRequest(r, http.MethodPost, "/vehicles", validVehicleBody)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestVehicleHandler_List(t *testing.T) {
	r := buildVehicleTestRouter()
	doRequest(r, http.MethodPost, "/vehicles", validVehicleBody)
	w := doRequest(r, http.MethodGet, "/vehicles", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestVehicleHandler_List_WithFilters(t *testing.T) {
	r := buildVehicleTestRouter()
	w := doRequest(r, http.MethodGet, "/vehicles?name=Civic&plate=ABC1234&brand=Honda&color=Prata&year=2020", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestVehicleHandler_List_InvalidYear(t *testing.T) {
	r := buildVehicleTestRouter()
	w := doRequest(r, http.MethodGet, "/vehicles?year=notanumber", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 (invalid year ignored), got %d", w.Code)
	}
}

func TestVehicleHandler_GetByID_Success(t *testing.T) {
	r := buildVehicleTestRouter()
	doRequest(r, http.MethodPost, "/vehicles", validVehicleBody)
	w := doRequest(r, http.MethodGet, "/vehicles/1", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestVehicleHandler_GetByID_NotFound(t *testing.T) {
	r := buildVehicleTestRouter()
	w := doRequest(r, http.MethodGet, "/vehicles/999", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestVehicleHandler_GetByID_InvalidID(t *testing.T) {
	r := buildVehicleTestRouter()
	w := doRequest(r, http.MethodGet, "/vehicles/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestVehicleHandler_Update_Success(t *testing.T) {
	r := buildVehicleTestRouter()
	doRequest(r, http.MethodPost, "/vehicles", validVehicleBody)
	body := `{"color":"Preto"}`
	w := doRequest(r, http.MethodPut, "/vehicles/1", body)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestVehicleHandler_Update_AllFields(t *testing.T) {
	r := buildVehicleTestRouter()
	doRequest(r, http.MethodPost, "/vehicles", validVehicleBody)
	body := `{"plate":"XYZ9876","name":"Fit","year":2021,"brand":"Honda","color":"Azul"}`
	w := doRequest(r, http.MethodPut, "/vehicles/1", body)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestVehicleHandler_Update_NotFound(t *testing.T) {
	r := buildVehicleTestRouter()
	w := doRequest(r, http.MethodPut, "/vehicles/999", `{"color":"Preto"}`)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestVehicleHandler_Update_InvalidID(t *testing.T) {
	r := buildVehicleTestRouter()
	w := doRequest(r, http.MethodPut, "/vehicles/abc", `{"color":"Preto"}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestVehicleHandler_Delete_Success(t *testing.T) {
	r := buildVehicleTestRouter()
	doRequest(r, http.MethodPost, "/vehicles", validVehicleBody)
	w := doRequest(r, http.MethodDelete, "/vehicles/1", "")
	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestVehicleHandler_Delete_NotFound(t *testing.T) {
	r := buildVehicleTestRouter()
	w := doRequest(r, http.MethodDelete, "/vehicles/999", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestVehicleHandler_Delete_InvalidID(t *testing.T) {
	r := buildVehicleTestRouter()
	w := doRequest(r, http.MethodDelete, "/vehicles/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
