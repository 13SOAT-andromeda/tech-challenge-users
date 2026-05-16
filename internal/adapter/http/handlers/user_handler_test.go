package handlers

import (
	"net/http"
	"strings"
	"testing"

	"tech-challenge-users/internal/application/services"

	"github.com/gin-gonic/gin"
)

func buildUserTestRouter() *gin.Engine {
	pr := newMockPersonRepo()
	ur := newMockUserRepo()
	er := newMockEmployeeRepo()
	tx := &mockTransactor{pr: pr, ur: ur, er: er}
	svc := services.NewUserService(tx, pr, ur, er)
	h := NewUserHandler(svc)
	ih := NewInternalUserHandler(svc)

	return setupRouter(func(g *gin.RouterGroup) {
		g.POST("/users", h.Create)
		g.GET("/users", h.List)
		g.GET("/users/:id", h.GetByID)
		g.PUT("/users/:id", h.Update)
		g.DELETE("/users/:id", h.Delete)
		g.GET("/internal/users/by-document", ih.GetByDocument)
	})
}

const validUserBody = `{"name":"Alice","email":"alice@example.com","password":"Alice@123","role":"attendant","document":"52998224725","position":"Attendant"}`

func TestUserHandler_Create_Success(t *testing.T) {
	r := buildUserTestRouter()
	w := doRequest(r, http.MethodPost, "/users", validUserBody)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUserHandler_Create_InvalidJSON(t *testing.T) {
	r := buildUserTestRouter()
	w := doRequest(r, http.MethodPost, "/users", `{invalid}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_Create_InvalidPassword(t *testing.T) {
	r := buildUserTestRouter()
	body := `{"name":"Alice","email":"alice@example.com","password":"weak","role":"attendant","document":"52998224725"}`
	w := doRequest(r, http.MethodPost, "/users", body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_Create_DuplicateEmail(t *testing.T) {
	r := buildUserTestRouter()
	doRequest(r, http.MethodPost, "/users", validUserBody)
	body2 := `{"name":"Bob","email":"alice@example.com","password":"Bob@12345","role":"attendant","document":"11144477735","position":"x"}`
	w := doRequest(r, http.MethodPost, "/users", body2)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestUserHandler_Create_DuplicateDocument(t *testing.T) {
	r := buildUserTestRouter()
	doRequest(r, http.MethodPost, "/users", validUserBody)
	body2 := `{"name":"Bob","email":"other@example.com","password":"Bob@12345","role":"attendant","document":"52998224725","position":"x"}`
	w := doRequest(r, http.MethodPost, "/users", body2)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestUserHandler_List(t *testing.T) {
	r := buildUserTestRouter()
	doRequest(r, http.MethodPost, "/users", validUserBody)
	w := doRequest(r, http.MethodGet, "/users", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUserHandler_List_WithFilters(t *testing.T) {
	r := buildUserTestRouter()
	w := doRequest(r, http.MethodGet, "/users?name=Alice&email=alice@example.com&contact=11999", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUserHandler_GetByID_Success(t *testing.T) {
	r := buildUserTestRouter()
	doRequest(r, http.MethodPost, "/users", validUserBody)
	w := doRequest(r, http.MethodGet, "/users/1", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUserHandler_GetByID_NotFound(t *testing.T) {
	r := buildUserTestRouter()
	w := doRequest(r, http.MethodGet, "/users/999", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUserHandler_GetByID_InvalidID(t *testing.T) {
	r := buildUserTestRouter()
	w := doRequest(r, http.MethodGet, "/users/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_Update_Success(t *testing.T) {
	r := buildUserTestRouter()
	doRequest(r, http.MethodPost, "/users", validUserBody)
	body := `{"name":"Alice Updated","contact":"11999999999"}`
	w := doRequest(r, http.MethodPut, "/users/1", body)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUserHandler_Update_WithAddress(t *testing.T) {
	r := buildUserTestRouter()
	doRequest(r, http.MethodPost, "/users", validUserBody)
	body := `{"name":"Alice","address":{"street":"Rua A","number":"10","city":"SP","country":"BR","zip_code":"01001000"}}`
	w := doRequest(r, http.MethodPut, "/users/1", body)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUserHandler_Update_PasswordRejected(t *testing.T) {
	r := buildUserTestRouter()
	doRequest(r, http.MethodPost, "/users", validUserBody)
	body := `{"password":"NewPass@123"}`
	w := doRequest(r, http.MethodPut, "/users/1", body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_Update_NotFound(t *testing.T) {
	r := buildUserTestRouter()
	w := doRequest(r, http.MethodPut, "/users/999", `{"name":"X"}`)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUserHandler_Update_InvalidID(t *testing.T) {
	r := buildUserTestRouter()
	w := doRequest(r, http.MethodPut, "/users/abc", `{"name":"X"}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_Delete_Success(t *testing.T) {
	r := buildUserTestRouter()
	doRequest(r, http.MethodPost, "/users", validUserBody)
	w := doRequest(r, http.MethodDelete, "/users/1", "")
	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestUserHandler_Delete_NotFound(t *testing.T) {
	r := buildUserTestRouter()
	w := doRequest(r, http.MethodDelete, "/users/999", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUserHandler_Delete_InvalidID(t *testing.T) {
	r := buildUserTestRouter()
	w := doRequest(r, http.MethodDelete, "/users/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- InternalUserHandler ---

func TestInternalUserHandler_GetByDocument_Success(t *testing.T) {
	r := buildUserTestRouter()
	doRequest(r, http.MethodPost, "/users", validUserBody)
	w := doRequest(r, http.MethodGet, "/internal/users/by-document?document=52998224725", "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestInternalUserHandler_GetByDocument_MissingParam(t *testing.T) {
	r := buildUserTestRouter()
	w := doRequest(r, http.MethodGet, "/internal/users/by-document", "")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestInternalUserHandler_GetByDocument_NotFound(t *testing.T) {
	r := buildUserTestRouter()
	w := doRequest(r, http.MethodGet, "/internal/users/by-document?document=52998224725", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// --- handleUserError coverage ---

func TestHandleUserError_DefaultBranch(t *testing.T) {
	r := buildUserTestRouter()
	body := `{"name":"Alice","email":"alice@example.com","password":"Alice@123","role":"invalid-role","document":"52998224725"}`
	w := doRequest(r, http.MethodPost, "/users", body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- toUserResponse with employee ---

func TestUserHandler_Create_CustomerHasNoEmployee(t *testing.T) {
	r := buildUserTestRouter()
	body := `{"name":"Bob","email":"bob@example.com","password":"Bob@12345","role":"customer","document":"52998224725"}`
	w := doRequest(r, http.MethodPost, "/users", body)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	if strings.Contains(w.Body.String(), `"employee":`) && !strings.Contains(w.Body.String(), `"employee":null`) {
		t.Error("expected no employee in response for customer role")
	}
}
