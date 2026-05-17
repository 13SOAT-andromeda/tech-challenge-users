package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"tech-challenge-users/internal/application/services"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	svc *services.EmployeeService
}

func NewEmployeeHandler(svc *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{svc: svc}
}

type employeeByIDResponse struct {
	ID       int64  `json:"id"`
	Position string `json:"position"`
	PersonID int64  `json:"person_id"`
}

func handleEmployeeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrEmployeeNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// GET /users/employees/:id
func (h *EmployeeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidID})
		return
	}

	employee, err := h.svc.GetByID(id)
	if err != nil {
		handleEmployeeError(c, err)
		return
	}
	c.JSON(http.StatusOK, employeeByIDResponse{
		ID:       employee.ID,
		Position: employee.Position,
		PersonID: employee.PersonID,
	})
}
