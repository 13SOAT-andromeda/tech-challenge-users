package e2e

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type employeeResp struct {
	ID       int64  `json:"id"`
	Position string `json:"position"`
	PersonID int64  `json:"person_id"`
}

func TestEmployee(t *testing.T) {
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

	t.Run("should get employee by ID successfully", func(t *testing.T) {
		// The admin user created on bootstrap is always employee ID=1
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/employees/1", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var employeeResponse employeeResp
		err = ParseBody(resp, &employeeResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, int64(1), employeeResponse.ID)
		assert.NotEmpty(t, employeeResponse.Position)
	})

	t.Run("should fail to get employee with invalid ID", func(t *testing.T) {
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/employees/invalid", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail to get non-existent employee", func(t *testing.T) {
		nonExistentID := strconv.FormatInt(999999, 10)
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/employees/"+nonExistentID, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
