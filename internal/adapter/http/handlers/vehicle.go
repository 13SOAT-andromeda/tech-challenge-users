package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/application/services"

	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	svc *services.VehicleService
}

func NewVehicleHandler(svc *services.VehicleService) *VehicleHandler {
	return &VehicleHandler{svc: svc}
}

type createVehicleRequest struct {
	Plate string `json:"plate" binding:"required"`
	Name  string `json:"name"  binding:"required"`
	Year  int    `json:"year"  binding:"required"`
	Brand string `json:"brand" binding:"required"`
	Color string `json:"color" binding:"required"`
}

type updateVehicleRequest struct {
	Plate *string `json:"plate"`
	Name  *string `json:"name"`
	Year  *int    `json:"year"`
	Brand *string `json:"brand"`
	Color *string `json:"color"`
}

type vehicleResponse struct {
	ID    int64  `json:"id"`
	Plate string `json:"plate"`
	Name  string `json:"name"`
	Year  int    `json:"year"`
	Brand string `json:"brand"`
	Color string `json:"color"`
}

func handleVehicleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrPlateAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrVehicleNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// POST /users/vehicles
func (h *VehicleHandler) Create(c *gin.Context) {
	var req createVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v, err := h.svc.CreateVehicle(services.CreateVehicleInput{
		Plate: req.Plate,
		Name:  req.Name,
		Year:  req.Year,
		Brand: req.Brand,
		Color: req.Color,
	})
	if err != nil {
		handleVehicleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, vehicleResponse{
		ID: v.ID, Plate: v.Plate, Name: v.Name,
		Year: v.Year, Brand: v.Brand, Color: v.Color,
	})
}

// GET /users/vehicles
func (h *VehicleHandler) List(c *gin.Context) {
	filters := ports.VehicleFilters{
		Name:  c.Query("name"),
		Plate: c.Query("plate"),
		Brand: c.Query("brand"),
		Color: c.Query("color"),
	}
	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			filters.Year = y
		}
	}

	vehicles, err := h.svc.ListVehicles(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	resp := make([]vehicleResponse, len(vehicles))
	for i, v := range vehicles {
		resp[i] = vehicleResponse{
			ID: v.ID, Plate: v.Plate, Name: v.Name,
			Year: v.Year, Brand: v.Brand, Color: v.Color,
		}
	}
	c.JSON(http.StatusOK, resp)
}

// GET /users/vehicles/:id
func (h *VehicleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	v, err := h.svc.GetByID(id)
	if err != nil {
		handleVehicleError(c, err)
		return
	}
	c.JSON(http.StatusOK, vehicleResponse{
		ID: v.ID, Plate: v.Plate, Name: v.Name,
		Year: v.Year, Brand: v.Brand, Color: v.Color,
	})
}

// PUT /users/vehicles/:id
func (h *VehicleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v, err := h.svc.UpdateVehicle(id, services.UpdateVehicleInput{
		Plate: req.Plate,
		Name:  req.Name,
		Year:  req.Year,
		Brand: req.Brand,
		Color: req.Color,
	})
	if err != nil {
		handleVehicleError(c, err)
		return
	}
	c.JSON(http.StatusOK, vehicleResponse{
		ID: v.ID, Plate: v.Plate, Name: v.Name,
		Year: v.Year, Brand: v.Brand, Color: v.Color,
	})
}

// DELETE /users/vehicles/:id
func (h *VehicleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.DeleteVehicle(id); err != nil {
		handleVehicleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
