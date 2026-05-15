package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/application/services"
	"tech-challenge-users/internal/domain"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *services.UserService
}

func NewUserHandler(svc *services.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// request / response types

type addressRequest struct {
	Street       string `json:"street"`
	Number       string `json:"number"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	Country      string `json:"country"`
	ZipCode      string `json:"zip_code"`
}

type createUserRequest struct {
	Name     string         `json:"name"     binding:"required"`
	Email    string         `json:"email"    binding:"required,email"`
	Password string         `json:"password" binding:"required"`
	Role     string         `json:"role"     binding:"required"`
	Document string         `json:"document" binding:"required"`
	Contact  string         `json:"contact"`
	Position string         `json:"position"`
	Address  addressRequest `json:"address"`
}

type updateUserRequest struct {
	Name     *string         `json:"name"`
	Email    *string         `json:"email"`
	Contact  *string         `json:"contact"`
	Position *string         `json:"position"`
	Address  *addressRequest `json:"address"`
	Password *string         `json:"password"` // always rejected with 400 if present
}

type addressResponse struct {
	Street       string `json:"street,omitempty"`
	Number       string `json:"number,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
	City         string `json:"city,omitempty"`
	Country      string `json:"country,omitempty"`
	ZipCode      string `json:"zip_code,omitempty"`
}

type personResponse struct {
	ID           int64           `json:"id"`
	Name         string          `json:"name"`
	Email        string          `json:"email"`
	Contact      string          `json:"contact,omitempty"`
	Document     string          `json:"document"`
	IsActive     bool            `json:"is_active"`
	Address      addressResponse `json:"address,omitempty"`
}

type employeeResponse struct {
	ID       int64  `json:"id"`
	Position string `json:"position"`
}

type userResponse struct {
	ID       int64             `json:"id"`
	Role     string            `json:"role"`
	Person   personResponse    `json:"person"`
	Employee *employeeResponse `json:"employee,omitempty"`
}

func toUserResponse(out services.UserOutput) userResponse {
	addr := addressResponse{
		Street:       out.Person.Address.Street,
		Number:       out.Person.Address.Number,
		Neighborhood: out.Person.Address.Neighborhood,
		City:         out.Person.Address.City,
		Country:      out.Person.Address.Country,
		ZipCode:      out.Person.Address.ZipCode,
	}
	resp := userResponse{
		ID:   out.User.ID,
		Role: out.User.Role,
		Person: personResponse{
			ID:       out.Person.ID,
			Name:     out.Person.Name,
			Email:    out.Person.Email,
			Contact:  out.Person.Contact,
			Document: out.Person.Document,
			IsActive: out.Person.IsActive,
			Address:  addr,
		},
	}
	if out.Employee != nil {
		resp.Employee = &employeeResponse{
			ID:       out.Employee.ID,
			Position: out.Employee.Position,
		}
	}
	return resp
}

func handleUserError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrEmailAlreadyExists),
		errors.Is(err, services.ErrDocumentAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// POST /users
func (h *UserHandler) Create(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.svc.CreateUser(services.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		Document: req.Document,
		Contact:  req.Contact,
		Position: req.Position,
		Address: domain.Address{
			Street:       req.Address.Street,
			Number:       req.Address.Number,
			Neighborhood: req.Address.Neighborhood,
			City:         req.Address.City,
			Country:      req.Address.Country,
			ZipCode:      req.Address.ZipCode,
		},
	})
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toUserResponse(*out))
}

// GET /users
func (h *UserHandler) List(c *gin.Context) {
	filters := ports.UserFilters{
		Name:    c.Query("name"),
		Email:   c.Query("email"),
		Contact: c.Query("contact"),
	}

	outputs, err := h.svc.ListUsers(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	resp := make([]userResponse, len(outputs))
	for i, out := range outputs {
		resp[i] = toUserResponse(out)
	}
	c.JSON(http.StatusOK, resp)
}

// GET /users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	out, err := h.svc.GetByID(id)
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusOK, toUserResponse(*out))
}

// PUT /users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Password != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password update is not allowed via this endpoint"})
		return
	}

	input := services.UpdateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Contact:  req.Contact,
		Position: req.Position,
	}
	if req.Address != nil {
		addr := domain.Address{
			Street:       req.Address.Street,
			Number:       req.Address.Number,
			Neighborhood: req.Address.Neighborhood,
			City:         req.Address.City,
			Country:      req.Address.Country,
			ZipCode:      req.Address.ZipCode,
		}
		input.Address = &addr
	}

	out, err := h.svc.UpdateUser(id, input)
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusOK, toUserResponse(*out))
}

// DELETE /users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.DeleteUser(id); err != nil {
		handleUserError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
