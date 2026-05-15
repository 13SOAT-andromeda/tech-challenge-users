package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"tech-challenge-users/internal/application/services"
	"tech-challenge-users/internal/domain"

	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	svc *services.CompanyService
}

func NewCompanyHandler(svc *services.CompanyService) *CompanyHandler {
	return &CompanyHandler{svc: svc}
}

type createCompanyRequest struct {
	Name     string         `json:"name"     binding:"required"`
	Email    string         `json:"email"    binding:"required,email"`
	Document string         `json:"document" binding:"required"`
	Contact  string         `json:"contact"  binding:"required"`
	Address  addressRequest `json:"address"`
}

type updateCompanyRequest struct {
	Name    *string         `json:"name"`
	Email   *string         `json:"email"`
	Contact *string         `json:"contact"`
	Address *addressRequest `json:"address"`
}

type companyResponse struct {
	ID       int64          `json:"id"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Document string         `json:"document"`
	Contact  string         `json:"contact"`
	Address  addressResponse `json:"address"`
}

func toCompanyResponse(c domain.Company) companyResponse {
	return companyResponse{
		ID:       c.ID,
		Name:     c.Name,
		Email:    c.Email,
		Document: c.Document,
		Contact:  c.Contact,
		Address: addressResponse{
			Street:       c.Address.Street,
			Number:       c.Address.Number,
			Neighborhood: c.Address.Neighborhood,
			City:         c.Address.City,
			Country:      c.Address.Country,
			ZipCode:      c.Address.ZipCode,
		},
	}
}

func handleCompanyError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrDocumentAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrCompanyNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// POST /users/companies
func (h *CompanyHandler) Create(c *gin.Context) {
	var req createCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company, err := h.svc.CreateCompany(services.CreateCompanyInput{
		Name:     req.Name,
		Email:    req.Email,
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
		handleCompanyError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toCompanyResponse(*company))
}

// GET /users/companies/:id
func (h *CompanyHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	company, err := h.svc.GetByID(id)
	if err != nil {
		handleCompanyError(c, err)
		return
	}
	c.JSON(http.StatusOK, toCompanyResponse(*company))
}

// PUT /users/companies/:id
func (h *CompanyHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := services.UpdateCompanyInput{
		Name:    req.Name,
		Email:   req.Email,
		Contact: req.Contact,
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

	company, err := h.svc.UpdateCompany(id, input)
	if err != nil {
		handleCompanyError(c, err)
		return
	}
	c.JSON(http.StatusOK, toCompanyResponse(*company))
}

// DELETE /users/companies/:id
func (h *CompanyHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.DeleteCompany(id); err != nil {
		handleCompanyError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
