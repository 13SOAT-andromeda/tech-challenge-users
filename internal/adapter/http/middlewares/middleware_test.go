package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

func newTestContext(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, nil)
	return c, w
}

// --- AuthRequired ---

func TestAuthRequired_WithUserID(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	c.Request.Header.Set("X-User-Id", "42")
	c.Request.Header.Set("X-User-Email", "user@example.com")
	c.Request.Header.Set("X-User-Role", "administrator")

	AuthRequired()(c)

	if w.Code == http.StatusUnauthorized {
		t.Error("expected request to be authorized when X-User-Id header is present")
	}
	if id, _ := c.Get(ContextUserID); id != "42" {
		t.Errorf("expected ContextUserID '42', got %v", id)
	}
	if email, _ := c.Get(ContextUserEmail); email != "user@example.com" {
		t.Errorf("expected ContextUserEmail to be set, got %v", email)
	}
	if role, _ := c.Get(ContextUserRole); role != "administrator" {
		t.Errorf("expected ContextUserRole to be set, got %v", role)
	}
}

func TestAuthRequired_MissingUserID(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")

	AuthRequired()(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// --- RoleRequired ---

func TestRoleRequired_AllowedRole(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	c.Set(ContextUserRole, "administrator")

	RoleRequired("administrator", "attendant")(c)

	if w.Code == http.StatusForbidden {
		t.Error("expected request to pass for allowed role")
	}
}

func TestRoleRequired_AnotherAllowedRole(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	c.Set(ContextUserRole, "attendant")

	RoleRequired("administrator", "attendant")(c)

	if w.Code == http.StatusForbidden {
		t.Error("expected request to pass for attendant role")
	}
}

func TestRoleRequired_ForbiddenRole(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	c.Set(ContextUserRole, "customer")

	RoleRequired("administrator")(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestRoleRequired_MissingRole(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")

	RoleRequired("administrator")(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}
