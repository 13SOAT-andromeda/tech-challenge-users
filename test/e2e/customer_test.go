package e2e

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type customerResp struct {
	ID     int64      `json:"id"`
	Type   string     `json:"type"`
	Person personResp `json:"person"`
}

type customerVehicleResp struct {
	ID         int64               `json:"id"`
	CustomerID int64               `json:"customer_id"`
	Vehicle    vehicleSummaryResp  `json:"vehicle"`
}

type vehicleSummaryResp struct {
	ID    int64  `json:"id"`
	Plate string `json:"plate"`
	Name  string `json:"name"`
	Year  int    `json:"year"`
	Brand string `json:"brand"`
	Color string `json:"color"`
}

type createCustomerReq struct {
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Type     string      `json:"type"`
	Document string      `json:"document"`
	Contact  string      `json:"contact"`
	Address  addressResp `json:"address"`
}

type updateCustomerReq struct {
	Name    *string      `json:"name,omitempty"`
	Contact *string      `json:"contact,omitempty"`
	Type    *string      `json:"type,omitempty"`
	Address *addressResp `json:"address,omitempty"`
}

func TestCustomer(t *testing.T) {
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

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	var createdCustomerID int64

	t.Run("should create a valid customer", func(t *testing.T) {
		req := createCustomerReq{
			Name:     "Maria Silva",
			Email:    "maria" + timestamp + "@test.com",
			Password: "Cassandra@123!",
			Type:     "pf",
			Document: GenerateValidCPF(time.Now().UnixNano()),
			Contact:  "11988887766",
			Address: addressResp{
				Street:       "Rua das Flores",
				Number:       "100",
				City:         "São Paulo",
				Neighborhood: "Jardins",
				Country:      "Brasil",
				ZipCode:      "01310100",
			},
		}

		payload, err := BuildBody(req)
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users/customers", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var customerResponse customerResp
		err = ParseBody(resp, &customerResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotZero(t, customerResponse.ID)
		assert.Equal(t, req.Type, customerResponse.Type)
		assert.Equal(t, req.Name, customerResponse.Person.Name)
		assert.Equal(t, req.Email, customerResponse.Person.Email)

		createdCustomerID = customerResponse.ID
	})

	t.Run("should fail to create customer with invalid request", func(t *testing.T) {
		req := createCustomerReq{
			Name:  "",
			Email: "invalid-email",
		}

		payload, err := BuildBody(req)
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users/customers", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail to create customer with duplicate email", func(t *testing.T) {
		req := createCustomerReq{
			Name:     "Duplicado",
			Email:    adminEmail,
			Password: "Cassandra@123!",
			Type:     "pf",
			Document: GenerateValidCPF(time.Now().UnixNano() + 2),
		}

		payload, err := BuildBody(req)
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users/customers", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("should get customer by ID successfully", func(t *testing.T) {
		customerId := strconv.FormatInt(createdCustomerID, 10)

		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/customers/"+customerId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var customerResponse customerResp
		err = ParseBody(resp, &customerResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, createdCustomerID, customerResponse.ID)
		assert.NotEmpty(t, customerResponse.Person.Name)
	})

	t.Run("should fail to get customer with invalid ID", func(t *testing.T) {
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/customers/invalid", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail to get non-existent customer", func(t *testing.T) {
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/customers/999999", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should list customers successfully", func(t *testing.T) {
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/customers", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var customersResponse []customerResp
		err = ParseBody(resp, &customersResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, customersResponse)
	})

	t.Run("should update customer successfully", func(t *testing.T) {
		updatedName := "Maria Silva Atualizada"
		updatedContact := "11977776655"
		updateInput := updateCustomerReq{
			Name:    &updatedName,
			Contact: &updatedContact,
		}

		payload, err := BuildBody(updateInput)
		require.NoError(t, err)

		customerId := strconv.FormatInt(createdCustomerID, 10)

		resp, err := NewIdentifiedReq("PUT", apiUrl+"/users/customers/"+customerId, payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var customerResponse customerResp
		err = ParseBody(resp, &customerResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, updatedName, customerResponse.Person.Name)
		assert.Equal(t, updatedContact, customerResponse.Person.Contact)
	})

	t.Run("should delete customer successfully", func(t *testing.T) {
		customerId := strconv.FormatInt(createdCustomerID, 10)

		resp, err := NewIdentifiedReq("DELETE", apiUrl+"/users/customers/"+customerId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("should verify deleted customer is not found", func(t *testing.T) {
		customerId := strconv.FormatInt(createdCustomerID, 10)

		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/customers/"+customerId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCustomerVehicle(t *testing.T) {
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

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	t.Run("should associate vehicle to customer and list", func(t *testing.T) {
		// Create customer
		customerReq := createCustomerReq{
			Name:     "Cliente Veículo",
			Email:    "clienteveiculo" + timestamp + "@test.com",
			Password: "Cassandra@123!",
			Type:     "pf",
			Document: GenerateValidCPF(time.Now().UnixNano()),
			Contact:  "11966665544",
		}
		payload, err := BuildBody(customerReq)
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users/customers", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var customer customerResp
		err = ParseBody(resp, &customer)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		customerID := customer.ID

		// Create vehicle
		vehicleReq := createVehicleReq{
			Plate: GenerateValidPlate(time.Now().UnixNano()),
			Name:  "Civic EX 2.0",
			Year:  2022,
			Brand: "Honda",
			Color: "Prata",
		}
		payload, err = BuildBody(vehicleReq)
		require.NoError(t, err)

		resp, err = NewIdentifiedReq("POST", apiUrl+"/users/vehicles", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var vehicle vehicleResp
		err = ParseBody(resp, &vehicle)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		vehicleID := vehicle.ID

		// Associate vehicle to customer
		customerIdStr := strconv.FormatInt(customerID, 10)
		vehicleIdStr := strconv.FormatInt(vehicleID, 10)

		resp, err = NewIdentifiedReq(
			"POST",
			apiUrl+"/users/customers/"+customerIdStr+"/vehicles/"+vehicleIdStr,
			nil, adminID, adminEmail, adminRole,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var cvResp customerVehicleResp
		err = ParseBody(resp, &cvResp)
		require.NoError(t, err)
		assert.Equal(t, customerID, cvResp.CustomerID)
		assert.Equal(t, vehicleID, cvResp.Vehicle.ID)

		// List customer vehicles
		resp, err = NewIdentifiedReq("GET", apiUrl+"/users/customers/"+customerIdStr+"/vehicles", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var cvList []customerVehicleResp
		err = ParseBody(resp, &cvList)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotEmpty(t, cvList, "customer should have at least one vehicle")

		var found bool
		for _, cv := range cvList {
			if cv.Vehicle.ID == vehicleID {
				found = true
				assert.Equal(t, customerID, cv.CustomerID)
				assert.Equal(t, vehicleReq.Plate, cv.Vehicle.Plate)
				break
			}
		}
		assert.True(t, found, "associated vehicle should appear in customer vehicles list")

		// Remove vehicle from customer
		resp, err = NewIdentifiedReq(
			"DELETE",
			apiUrl+"/users/customers/"+customerIdStr+"/vehicles/"+vehicleIdStr,
			nil, adminID, adminEmail, adminRole,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("should fail to associate vehicle to non-existent customer", func(t *testing.T) {
		resp, err := NewIdentifiedReq(
			"POST",
			apiUrl+"/users/customers/999999/vehicles/999999",
			nil, adminID, adminEmail, adminRole,
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
