package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"tech-challenge-users/internal/application/usecases"

	"github.com/gin-gonic/gin"
)

const errInvalidCustomerID = "invalid customer id"

type CustomerVehicleHandler struct {
	uc *usecases.CustomerUseCase
}

func NewCustomerVehicleHandler(uc *usecases.CustomerUseCase) *CustomerVehicleHandler {
	return &CustomerVehicleHandler{uc: uc}
}

type vehicleSummaryResponse struct {
	ID    int64  `json:"id"`
	Plate string `json:"plate"`
	Name  string `json:"name"`
	Year  int    `json:"year"`
	Brand string `json:"brand"`
	Color string `json:"color"`
}

type customerVehicleResponse struct {
	ID         int64                  `json:"id"`
	CustomerID int64                  `json:"customer_id"`
	Vehicle    vehicleSummaryResponse `json:"vehicle"`
}

func handleCVError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecases.ErrCustomerNotFound), errors.Is(err, usecases.ErrVehicleNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, usecases.ErrAssociationAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, usecases.ErrAssociationNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

// GET /users/customers/:id/vehicles
func (h *CustomerVehicleHandler) List(c *gin.Context) {
	customerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidCustomerID})
		return
	}

	cvs, err := h.uc.GetCustomerVehicles(customerID)
	if err != nil {
		handleCVError(c, err)
		return
	}

	resp := make([]customerVehicleResponse, len(cvs))
	for i, cv := range cvs {
		resp[i] = customerVehicleResponse{
			ID:         cv.ID,
			CustomerID: cv.CustomerID,
			Vehicle: vehicleSummaryResponse{
				ID:    cv.Vehicle.ID,
				Plate: cv.Vehicle.Plate,
				Name:  cv.Vehicle.Name,
				Year:  cv.Vehicle.Year,
				Brand: cv.Vehicle.Brand,
				Color: cv.Vehicle.Color,
			},
		}
	}
	c.JSON(http.StatusOK, resp)
}

// POST /users/customers/:id/vehicles/:vehicleId
func (h *CustomerVehicleHandler) Add(c *gin.Context) {
	customerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidCustomerID})
		return
	}

	vehicleID, err := strconv.ParseInt(c.Param("vehicleId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle id"})
		return
	}

	cv, err := h.uc.AddVehicle(customerID, vehicleID)
	if err != nil {
		handleCVError(c, err)
		return
	}

	c.JSON(http.StatusCreated, customerVehicleResponse{
		ID:         cv.ID,
		CustomerID: cv.CustomerID,
		Vehicle: vehicleSummaryResponse{
			ID:    cv.Vehicle.ID,
			Plate: cv.Vehicle.Plate,
			Name:  cv.Vehicle.Name,
			Year:  cv.Vehicle.Year,
			Brand: cv.Vehicle.Brand,
			Color: cv.Vehicle.Color,
		},
	})
}

// DELETE /users/customers/:id/vehicles/:vehicleId
func (h *CustomerVehicleHandler) Remove(c *gin.Context) {
	customerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidCustomerID})
		return
	}

	vehicleID, err := strconv.ParseInt(c.Param("vehicleId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle id"})
		return
	}

	if err := h.uc.RemoveVehicle(customerID, vehicleID); err != nil {
		handleCVError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
