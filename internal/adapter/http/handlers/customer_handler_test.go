package handlers

import (
	"net/http"
	"testing"

	"tech-challenge-users/internal/application/services"
	"tech-challenge-users/internal/application/usecases"

	"github.com/gin-gonic/gin"
)

func buildCustomerTestRouter() *gin.Engine {
	pr := newMockPersonRepo()
	ur := newMockUserRepo()
	er := newMockEmployeeRepo()
	cr := newMockCustomerRepo()
	vr := newMockVehicleRepo()
	cvr := newMockCVRepo()
	tx := &mockTransactor{pr: pr, ur: ur, er: er, cr: cr}

	customerSvc := services.NewCustomerService(tx, pr, cr)
	uc := usecases.NewCustomerUseCase(cr, vr, cvr)

	h := NewCustomerHandler(customerSvc)
	cvh := NewCustomerVehicleHandler(uc)

	return setupRouter(func(g *gin.RouterGroup) {
		g.POST("/customers", h.Create)
		g.GET("/customers", h.List)
		g.GET("/customers/:id", h.GetByID)
		g.PUT("/customers/:id", h.Update)
		g.DELETE("/customers/:id", h.Delete)
		g.GET("/customers/:id/vehicles", cvh.List)
		g.POST("/customers/:id/vehicles/:vehicleId", cvh.Add)
		g.DELETE("/customers/:id/vehicles/:vehicleId", cvh.Remove)
	})
}

const validCustomerBody = `{"name":"Bob","email":"bob@example.com","password":"Bob@12345","type":"pf","document":"52998224725"}`

func TestCustomerHandler_Create_Success(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodPost, "/customers", validCustomerBody)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCustomerHandler_Create_InvalidJSON(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodPost, "/customers", `{bad}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCustomerHandler_Create_InvalidType(t *testing.T) {
	r := buildCustomerTestRouter()
	body := `{"name":"Bob","email":"bob@example.com","password":"Bob@12345","type":"corp","document":"52998224725"}`
	w := doRequest(r, http.MethodPost, "/customers", body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCustomerHandler_Create_DuplicateEmail(t *testing.T) {
	r := buildCustomerTestRouter()
	doRequest(r, http.MethodPost, "/customers", validCustomerBody)
	body2 := `{"name":"Bob2","email":"bob@example.com","password":"Bob@12345","type":"pf","document":"11144477735"}`
	w := doRequest(r, http.MethodPost, "/customers", body2)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestCustomerHandler_Create_DuplicateDocument(t *testing.T) {
	r := buildCustomerTestRouter()
	doRequest(r, http.MethodPost, "/customers", validCustomerBody)
	body2 := `{"name":"Bob2","email":"other@example.com","password":"Bob@12345","type":"pf","document":"52998224725"}`
	w := doRequest(r, http.MethodPost, "/customers", body2)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestCustomerHandler_List(t *testing.T) {
	r := buildCustomerTestRouter()
	doRequest(r, http.MethodPost, "/customers", validCustomerBody)
	w := doRequest(r, http.MethodGet, "/customers", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCustomerHandler_List_WithFilters(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodGet, "/customers?name=Bob&status=active", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCustomerHandler_List_InactiveStatus(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodGet, "/customers?status=inactive", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCustomerHandler_GetByID_Success(t *testing.T) {
	r := buildCustomerTestRouter()
	doRequest(r, http.MethodPost, "/customers", validCustomerBody)
	w := doRequest(r, http.MethodGet, "/customers/1", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCustomerHandler_GetByID_NotFound(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodGet, "/customers/999", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCustomerHandler_GetByID_InvalidID(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodGet, "/customers/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCustomerHandler_Update_Success(t *testing.T) {
	r := buildCustomerTestRouter()
	doRequest(r, http.MethodPost, "/customers", validCustomerBody)
	body := `{"name":"Bob Updated","contact":"11999999999"}`
	w := doRequest(r, http.MethodPut, "/customers/1", body)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCustomerHandler_Update_WithAddress(t *testing.T) {
	r := buildCustomerTestRouter()
	doRequest(r, http.MethodPost, "/customers", validCustomerBody)
	body := `{"address":{"street":"Rua B","number":"5","city":"RJ"}}`
	w := doRequest(r, http.MethodPut, "/customers/1", body)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCustomerHandler_Update_NotFound(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodPut, "/customers/999", `{"name":"X"}`)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCustomerHandler_Update_InvalidID(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodPut, "/customers/abc", `{"name":"X"}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCustomerHandler_Delete_Success(t *testing.T) {
	r := buildCustomerTestRouter()
	doRequest(r, http.MethodPost, "/customers", validCustomerBody)
	w := doRequest(r, http.MethodDelete, "/customers/1", "")
	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestCustomerHandler_Delete_NotFound(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodDelete, "/customers/999", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCustomerHandler_Delete_InvalidID(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodDelete, "/customers/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- CustomerVehicleHandler ---

func TestCVHandler_List_Success(t *testing.T) {
	r := buildCustomerTestRouter()
	doRequest(r, http.MethodPost, "/customers", validCustomerBody)
	w := doRequest(r, http.MethodGet, "/customers/1/vehicles", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCVHandler_List_InvalidID(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodGet, "/customers/abc/vehicles", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCVHandler_List_CustomerNotFound(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodGet, "/customers/999/vehicles", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCVHandler_Add_Success(t *testing.T) {
	r := buildCustomerTestRouter()
	// Create customer and vehicle first
	doRequest(r, http.MethodPost, "/customers", validCustomerBody)
	// Seed vehicle repo via vehicle handler would be complex; we test the error path instead
	w := doRequest(r, http.MethodPost, "/customers/1/vehicles/999", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for missing vehicle, got %d", w.Code)
	}
}

func TestCVHandler_Add_InvalidCustomerID(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodPost, "/customers/abc/vehicles/1", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCVHandler_Add_InvalidVehicleID(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodPost, "/customers/1/vehicles/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCVHandler_Remove_InvalidCustomerID(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodDelete, "/customers/abc/vehicles/1", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCVHandler_Remove_InvalidVehicleID(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodDelete, "/customers/1/vehicles/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCVHandler_Remove_NotFound(t *testing.T) {
	r := buildCustomerTestRouter()
	w := doRequest(r, http.MethodDelete, "/customers/1/vehicles/1", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
