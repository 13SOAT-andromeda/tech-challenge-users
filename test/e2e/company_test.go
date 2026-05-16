package e2e

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type companyResp struct {
	ID       int64       `json:"id"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Document string      `json:"document"`
	Contact  string      `json:"contact"`
	Address  addressResp `json:"address"`
}

type createCompanyReq struct {
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Document string      `json:"document"`
	Contact  string      `json:"contact"`
	Address  addressResp `json:"address"`
}

type updateCompanyReq struct {
	Name    *string      `json:"name,omitempty"`
	Email   *string      `json:"email,omitempty"`
	Contact *string      `json:"contact,omitempty"`
	Address *addressResp `json:"address,omitempty"`
}

func TestCompany(t *testing.T) {
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
	var createdCompanyID int64

	t.Run("should create a valid company", func(t *testing.T) {
		req := createCompanyReq{
			Name:     "Oficina Teste",
			Email:    "oficina" + timestamp + "@test.com",
			Document: GenerateValidCNPJ(time.Now().UnixNano()),
			Contact:  "11933334444",
			Address: addressResp{
				Street:       "Av. Paulista",
				Number:       "1000",
				City:         "São Paulo",
				Neighborhood: "Bela Vista",
				Country:      "Brasil",
				ZipCode:      "01310100",
			},
		}

		payload, err := BuildBody(req)
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users/companies", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var companyResponse companyResp
		err = ParseBody(resp, &companyResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotZero(t, companyResponse.ID)
		assert.Equal(t, req.Name, companyResponse.Name)
		assert.Equal(t, req.Email, companyResponse.Email)
		assert.Equal(t, req.Document, companyResponse.Document)
		assert.Equal(t, req.Contact, companyResponse.Contact)

		createdCompanyID = companyResponse.ID
	})

	t.Run("should fail to create company with invalid request", func(t *testing.T) {
		req := createCompanyReq{
			Name:  "",
			Email: "invalid-email",
		}

		payload, err := BuildBody(req)
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users/companies", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail to create company without admin role", func(t *testing.T) {
		req := createCompanyReq{
			Name:     "Empresa Não Autorizada",
			Email:    "nauto" + timestamp + "@test.com",
			Document: "99999999000199",
			Contact:  "11900000000",
		}

		payload, err := BuildBody(req)
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users/companies", payload, adminID, adminEmail, "attendant")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should get company by ID successfully", func(t *testing.T) {
		companyId := strconv.FormatInt(createdCompanyID, 10)

		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/companies/"+companyId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var companyResponse companyResp
		err = ParseBody(resp, &companyResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, createdCompanyID, companyResponse.ID)
		assert.NotEmpty(t, companyResponse.Name)
	})

	t.Run("should fail to get company with invalid ID", func(t *testing.T) {
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/companies/invalid", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail to get non-existent company", func(t *testing.T) {
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/companies/999999", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should update company successfully", func(t *testing.T) {
		updatedName := "Oficina Teste Atualizada"
		updatedContact := "11944445555"
		updateInput := updateCompanyReq{
			Name:    &updatedName,
			Contact: &updatedContact,
		}

		payload, err := BuildBody(updateInput)
		require.NoError(t, err)

		companyId := strconv.FormatInt(createdCompanyID, 10)

		resp, err := NewIdentifiedReq("PUT", apiUrl+"/users/companies/"+companyId, payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var companyResponse companyResp
		err = ParseBody(resp, &companyResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, updatedName, companyResponse.Name)
		assert.Equal(t, updatedContact, companyResponse.Contact)
	})

	t.Run("should delete company successfully", func(t *testing.T) {
		companyId := strconv.FormatInt(createdCompanyID, 10)

		resp, err := NewIdentifiedReq("DELETE", apiUrl+"/users/companies/"+companyId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("should verify deleted company is not found", func(t *testing.T) {
		companyId := strconv.FormatInt(createdCompanyID, 10)

		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/companies/"+companyId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
