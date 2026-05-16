package e2e

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type addressResp struct {
	Street       string `json:"street"`
	Number       string `json:"number"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	Country      string `json:"country"`
	ZipCode      string `json:"zip_code"`
}

type personResp struct {
	ID       int64       `json:"id"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Contact  string      `json:"contact"`
	Document string      `json:"document"`
	IsActive bool        `json:"is_active"`
	Address  addressResp `json:"address"`
}

type userResp struct {
	ID     int64      `json:"id"`
	Role   string     `json:"role"`
	Person personResp `json:"person"`
}

type createUserReq struct {
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Role     string      `json:"role"`
	Document string      `json:"document"`
	Contact  string      `json:"contact"`
	Position string      `json:"position"`
	Address  addressResp `json:"address"`
}

type updateUserReq struct {
	Name     *string      `json:"name,omitempty"`
	Contact  *string      `json:"contact,omitempty"`
	Position *string      `json:"position,omitempty"`
	Address  *addressResp `json:"address,omitempty"`
}

func TestUser(t *testing.T) {
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

	var createdUserID int64

	t.Run("should create a valid user", func(t *testing.T) {
		req := createUserReq{
			Name:     "João Marcos",
			Email:    "marcosjoao" + strconv.FormatInt(time.Now().Unix(), 10) + "@gmail.com",
			Role:     "attendant",
			Document: GenerateValidCPF(time.Now().UnixNano()),
			Contact:  "11979664877",
			Password: "Cassandra@123!",
			Position: "Atendente",
			Address: addressResp{
				Street:       "Rua Diamantina",
				Number:       "430",
				City:         "São Paulo",
				Neighborhood: "Vila Maria",
				Country:      "Brasil",
				ZipCode:      "02170150",
			},
		}

		payload, err := BuildBody(req)
		require.NoError(t, err, "failed to build payload")

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err, "failed to make create user request")
		defer resp.Body.Close()

		var userResponse userResp
		err = ParseBody(resp, &userResponse)
		require.NoError(t, err, "failed to parse create user response")

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotZero(t, userResponse.ID)
		assert.Equal(t, req.Role, userResponse.Role)
		assert.Equal(t, req.Name, userResponse.Person.Name)
		assert.Equal(t, req.Email, userResponse.Person.Email)
		assert.Equal(t, req.Contact, userResponse.Person.Contact)

		createdUserID = userResponse.ID
	})

	t.Run("should fail to create user with invalid request", func(t *testing.T) {
		req := createUserReq{
			Name:  "",
			Email: "invalid-email",
		}

		payload, err := BuildBody(req)
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail to create user with duplicate email", func(t *testing.T) {
		req := createUserReq{
			Name:     "João Duplicado",
			Email:    adminEmail,
			Role:     "attendant",
			Document: GenerateValidCPF(time.Now().UnixNano() + 1),
			Contact:  "11979664877",
			Password: "Cassandra@123!",
		}

		payload, err := BuildBody(req)
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("should get user by ID successfully", func(t *testing.T) {
		userId := strconv.FormatInt(createdUserID, 10)

		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/"+userId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err, "failed to make get user request")
		defer resp.Body.Close()

		var userResponse userResp
		err = ParseBody(resp, &userResponse)
		require.NoError(t, err, "failed to parse get user response")

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, createdUserID, userResponse.ID)
		assert.NotEmpty(t, userResponse.Person.Name)
		assert.NotEmpty(t, userResponse.Person.Email)
	})

	t.Run("should fail to get user with invalid ID", func(t *testing.T) {
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/invalid", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail to get non-existent user", func(t *testing.T) {
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/999999", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should search users successfully", func(t *testing.T) {
		params := url.Values{}
		params.Add("name", "João")
		searchUrl := apiUrl + "/users?" + params.Encode()

		resp, err := NewIdentifiedReq("GET", searchUrl, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var usersResponse []userResp
		err = ParseBody(resp, &usersResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, usersResponse)
	})

	t.Run("should search users without parameters", func(t *testing.T) {
		resp, err := NewIdentifiedReq("GET", apiUrl+"/users", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var usersResponse []userResp
		err = ParseBody(resp, &usersResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, usersResponse)
	})

	t.Run("should update user successfully", func(t *testing.T) {
		updatedName := "João Marcos Atualizado"
		updatedContact := "11988888888"
		updateInput := updateUserReq{
			Name:    &updatedName,
			Contact: &updatedContact,
			Address: &addressResp{
				Street:       "Rua Atualizada",
				Number:       "500",
				City:         "Rio de Janeiro",
				Neighborhood: "Centro",
				Country:      "Brasil",
				ZipCode:      "20000000",
			},
		}

		payload, err := BuildBody(updateInput)
		require.NoError(t, err)

		userId := strconv.FormatInt(createdUserID, 10)

		resp, err := NewIdentifiedReq("PUT", apiUrl+"/users/"+userId, payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		var userResponse userResp
		err = ParseBody(resp, &userResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, updatedName, userResponse.Person.Name)
		assert.Equal(t, updatedContact, userResponse.Person.Contact)
	})

	t.Run("should fail to update user with invalid ID", func(t *testing.T) {
		updatedName := "Test"
		payload, err := BuildBody(updateUserReq{Name: &updatedName})
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("PUT", apiUrl+"/users/invalid", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail to update non-existent user", func(t *testing.T) {
		updatedName := "Test User"
		payload, err := BuildBody(updateUserReq{Name: &updatedName})
		require.NoError(t, err)

		resp, err := NewIdentifiedReq("PUT", apiUrl+"/users/999999", payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should reject password update", func(t *testing.T) {
		newPassword := "NewPassword@123!"
		payload, err := BuildBody(map[string]interface{}{"password": newPassword})
		require.NoError(t, err)

		userId := strconv.FormatInt(createdUserID, 10)

		resp, err := NewIdentifiedReq("PUT", apiUrl+"/users/"+userId, payload, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should delete user successfully", func(t *testing.T) {
		userId := strconv.FormatInt(createdUserID, 10)

		resp, err := NewIdentifiedReq("DELETE", apiUrl+"/users/"+userId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("should fail to delete user with invalid ID", func(t *testing.T) {
		resp, err := NewIdentifiedReq("DELETE", apiUrl+"/users/invalid", nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should verify deleted user is not found", func(t *testing.T) {
		userId := strconv.FormatInt(createdUserID, 10)

		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/"+userId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
