package handlers

import (
	"net/http"

	"tech-challenge-users/internal/application/services"

	"github.com/gin-gonic/gin"
)

type InternalUserHandler struct {
	svc *services.UserService
}

func NewInternalUserHandler(svc *services.UserService) *InternalUserHandler {
	return &InternalUserHandler{svc: svc}
}

// GET /internal/users/by-document?document=<cpf>
func (h *InternalUserHandler) GetByDocument(c *gin.Context) {
	document := c.Query("document")
	if document == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document is required"})
		return
	}

	out, err := h.svc.GetByDocument(document)
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"person": gin.H{
			"id":           out.Person.ID,
			"name":         out.Person.Name,
			"email":        out.Person.Email,
			"contact":      out.Person.Contact,
			"document":     out.Person.Document,
			"is_active":    out.Person.IsActive,
			"address":      out.Person.Address.Street,
			"address_number": out.Person.Address.Number,
			"neighborhood": out.Person.Address.Neighborhood,
			"city":         out.Person.Address.City,
			"country":      out.Person.Address.Country,
			"zip_code":     out.Person.Address.ZipCode,
			"created_at":   out.Person.CreatedAt,
			"updated_at":   out.Person.UpdatedAt,
		},
		"user": gin.H{
			"id":        out.User.ID,
			"password":  out.User.Password,
			"role":      out.User.Role,
			"person_id": out.User.PersonID,
		},
	})
}
