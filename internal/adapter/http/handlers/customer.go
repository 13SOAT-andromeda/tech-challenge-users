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

const errInvalidID = "invalid id"

type CustomerHandler struct {
	svc *services.CustomerService
}

func NewCustomerHandler(svc *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

type createCustomerRequest struct {
	Name     string         `json:"name"     binding:"required"`
	Email    string         `json:"email"    binding:"required,email"`
	Password string         `json:"password" binding:"required"`
	Type     string         `json:"type"     binding:"required"`
	Document string         `json:"document" binding:"required"`
	Contact  string         `json:"contact"`
	Address  addressRequest `json:"address"`
}

type updateCustomerRequest struct {
	Name    *string         `json:"name"`
	Email   *string         `json:"email"`
	Contact *string         `json:"contact"`
	Address *addressRequest `json:"address"`
	Type    *string         `json:"type"`
}

type customerResponse struct {
	ID     int64          `json:"id"`
	Type   string         `json:"type"`
	Person personResponse `json:"person"`
}

func toCustomerResponse(out services.CustomerOutput) customerResponse {
	return customerResponse{
		ID:   out.Customer.ID,
		Type: out.Customer.Type,
		Person: personResponse{
			ID:       out.Person.ID,
			Name:     out.Person.Name,
			Email:    out.Person.Email,
			Contact:  out.Person.Contact,
			Document: out.Person.Document,
			IsActive: out.Person.IsActive,
			Address: addressResponse{
				Street:       out.Person.Address.Street,
				Number:       out.Person.Address.Number,
				Neighborhood: out.Person.Address.Neighborhood,
				City:         out.Person.Address.City,
				Country:      out.Person.Address.Country,
				ZipCode:      out.Person.Address.ZipCode,
			},
		},
	}
}

func handleCustomerError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrEmailAlreadyExists),
		errors.Is(err, services.ErrDocumentAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrCustomerNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// POST /users/customers
func (h *CustomerHandler) Create(c *gin.Context) {
	var req createCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.svc.CreateCustomer(services.CreateCustomerInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Type:     req.Type,
		Document: req.Document,
		Contact:  req.Contact,
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
		handleCustomerError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toCustomerResponse(*out))
}

// GET /users/customers
func (h *CustomerHandler) List(c *gin.Context) {
	filters := ports.CustomerFilters{
		Document: c.Query("document"),
		Name:     c.Query("name"),
		Email:    c.Query("email"),
	}
	if status := c.Query("status"); status != "" {
		active := status == "active"
		filters.IsActive = &active
	}

	outputs, err := h.svc.ListCustomers(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	resp := make([]customerResponse, len(outputs))
	for i, out := range outputs {
		resp[i] = toCustomerResponse(out)
	}
	c.JSON(http.StatusOK, resp)
}

// GET /users/customers/:id
func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidID})
		return
	}

	out, err := h.svc.GetByID(id)
	if err != nil {
		handleCustomerError(c, err)
		return
	}
	c.JSON(http.StatusOK, toCustomerResponse(*out))
}

// PUT /users/customers/:id
func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidID})
		return
	}

	var req updateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := services.UpdateCustomerInput{
		Name:    req.Name,
		Email:   req.Email,
		Contact: req.Contact,
		Type:    req.Type,
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

	out, err := h.svc.UpdateCustomer(id, input)
	if err != nil {
		handleCustomerError(c, err)
		return
	}
	c.JSON(http.StatusOK, toCustomerResponse(*out))
}

// DELETE /users/customers/:id
func (h *CustomerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidID})
		return
	}

	if err := h.svc.DeleteCustomer(id); err != nil {
		handleCustomerError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
