package e2e

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type vehicleResp struct {
	ID    int64  `json:"id"`
	Plate string `json:"plate"`
	Name  string `json:"name"`
	Year  int    `json:"year"`
	Brand string `json:"brand"`
	Color string `json:"color"`
}

type createVehicleReq struct {
	Plate string `json:"plate"`
	Name  string `json:"name"`
	Year  int    `json:"year"`
	Brand string `json:"brand"`
	Color string `json:"color"`
}

type updateVehicleReq struct {
	Plate *string `json:"plate,omitempty"`
	Name  *string `json:"name,omitempty"`
	Year  *int    `json:"year,omitempty"`
	Brand *string `json:"brand,omitempty"`
	Color *string `json:"color,omitempty"`
}

func TestVehicle(t *testing.T) {
	cfg, err := SetupTest()
	if err != nil {
		t.Skipf("Error on setup config variables: %v", err)
	}

	apiUrl := GetApiUrl(*cfg)

	healthResp, err := http.Get(apiUrl + "/health")
	if err != nil {
		t.Skipf("Application not running at %s, skipping E2E tests: %v", apiUrl, err)
		return
	}
	healthResp.Body.Close()

	adminEmail := cfg.AdminEmail
	adminRole := "administrator"
	adminID := "1"

	var createdVehicleID int64

	t.Run("should create a valid vehicle", func(t *testing.T) {
		req := createVehicleReq{
			Plate: GenerateValidPlate(time.Now().UnixNano()),
			Name:  "Corolla XEI 2.0",
			Year:  2021,
			Brand: "Toyota",
			Color: "Branco",
		}

		payload, err := BuildBody(req)
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users/vehicles", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var vehicleResponse vehicleResp
		err = ParseBody(resp, &vehicleResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotZero(t, vehicleResponse.ID)
		assert.Equal(t, req.Plate, vehicleResponse.Plate)
		assert.Equal(t, req.Name, vehicleResponse.Name)
		assert.Equal(t, req.Year, vehicleResponse.Year)
		assert.Equal(t, req.Brand, vehicleResponse.Brand)
		assert.Equal(t, req.Color, vehicleResponse.Color)

		createdVehicleID = vehicleResponse.ID
	})

	t.Run("should fail to create vehicle with invalid request", func(t *testing.T) {
		req := createVehicleReq{
			Name: "",
		}

		payload, err := BuildBody(req)
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users/vehicles", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should get vehicle by ID successfully", func(t *testing.T) {
		vehicleId := strconv.FormatInt(createdVehicleID, 10)

		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/vehicles/"+vehicleId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var vehicleResponse vehicleResp
		err = ParseBody(resp, &vehicleResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, createdVehicleID, vehicleResponse.ID)
		assert.NotEmpty(t, vehicleResponse.Plate)
	})

	t.Run("should fail to get vehicle with invalid ID", func(t *testing.T) {
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/vehicles/invalid", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail to get non-existent vehicle", func(t *testing.T) {
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/vehicles/999999", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should list vehicles successfully", func(t *testing.T) {
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/vehicles", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var vehiclesResponse []vehicleResp
		err = ParseBody(resp, &vehiclesResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, vehiclesResponse)
	})

	t.Run("should update vehicle successfully", func(t *testing.T) {
		updatedColor := "Preto"
		updatedName := "Corolla XEI 2.0 Flex"
		updateInput := updateVehicleReq{
			Color: &updatedColor,
			Name:  &updatedName,
		}

		payload, err := BuildBody(updateInput)
		require.NoError(t, err)

		vehicleId := strconv.FormatInt(createdVehicleID, 10)

		resp, err := NewIdentifiedReq("PUT", apiUrl+"/users/vehicles/"+vehicleId, payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var vehicleResponse vehicleResp
		err = ParseBody(resp, &vehicleResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, updatedColor, vehicleResponse.Color)
		assert.Equal(t, updatedName, vehicleResponse.Name)
	})

	t.Run("should fail to update vehicle with invalid ID", func(t *testing.T) {
		updatedColor := "Azul"
		payload, err := BuildBody(updateVehicleReq{Color: &updatedColor})
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("PUT", apiUrl+"/users/vehicles/invalid", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail to create vehicle with duplicate plate", func(t *testing.T) {
		// Get existing plate first
		vehicleId := strconv.FormatInt(createdVehicleID, 10)
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/vehicles/"+vehicleId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var existingVehicle vehicleResp
		err = ParseBody(resp, &existingVehicle)
		require.NoError(t, err)

		req := createVehicleReq{
			Plate: existingVehicle.Plate,
			Name:  "Outro Carro",
			Year:  2020,
			Brand: "Ford",
			Color: "Vermelho",
		}

		payload, err := BuildBody(req)
		require.NoError(t, err)

		resp, err = NewIdentifiedReq("POST", apiUrl+"/users/vehicles", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("should delete vehicle successfully", func(t *testing.T) {
		vehicleId := strconv.FormatInt(createdVehicleID, 10)

		resp, err := NewIdentifiedReq("DELETE", apiUrl+"/users/vehicles/"+vehicleId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("should verify deleted vehicle is not found", func(t *testing.T) {
		vehicleId := strconv.FormatInt(createdVehicleID, 10)

		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/vehicles/"+vehicleId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
