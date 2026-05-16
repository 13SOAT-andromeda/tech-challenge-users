package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
)

type testConfig struct {
	APIURL     string
	AdminEmail string
}

// SetupTest loads test configuration from environment variables.
// API_URL defaults to http://localhost:8081
// ADMIN_EMAIL defaults to admin@admin.com
func SetupTest() (*testConfig, error) {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8081"
	}
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@garage.com"
	}
	return &testConfig{
		APIURL:     apiURL,
		AdminEmail: adminEmail,
	}, nil
}

func GetApiUrl(cfg testConfig) string {
	return cfg.APIURL
}

// NewUnauthenticatedReq creates an HTTP request without identity headers.
func NewUnauthenticatedReq(method, url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	if body == nil {
		body = http.NoBody
	}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	return client.Do(request)
}

// NewIdentifiedReq creates a request with Lambda Authorizer identity headers.
// userID, email and role are the values the Lambda would inject after JWT validation.
func NewIdentifiedReq(method, url string, body io.Reader, userID, email, role string) (*http.Response, error) {
	client := &http.Client{}
	if body == nil {
		body = http.NoBody
	}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-User-Id", userID)
	request.Header.Add("X-User-Email", email)
	request.Header.Add("X-User-Role", role)
	return client.Do(request)
}

func ParseBody[T any](resp *http.Response, data *T) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.Unmarshal(bodyBytes, data)
}

func BuildBody[T any](data T) (io.Reader, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonData), nil
}

func GenerateValidCPF(seed int64) string {
	base := seed % 1000000000
	if base < 100000000 {
		base += 100000000
	}

	baseStr := strconv.FormatInt(base, 10)

	firstChar := baseStr[0]
	allSame := true
	for i := 1; i < len(baseStr); i++ {
		if baseStr[i] != firstChar {
			allSame = false
			break
		}
	}

	if allSame {
		lastDigit, _ := strconv.Atoi(string(baseStr[8]))
		lastDigit = (lastDigit + 1) % 10
		baseStr = baseStr[:8] + strconv.Itoa(lastDigit)
	}

	sum := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(baseStr[i]))
		sum += digit * (10 - i)
	}

	firstDigit := (sum * 10) % 11
	if firstDigit == 10 {
		firstDigit = 0
	}

	baseStr += strconv.Itoa(firstDigit)
	sum = 0
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(baseStr[i]))
		sum += digit * (11 - i)
	}

	secondDigit := (sum * 10) % 11
	if secondDigit == 10 {
		secondDigit = 0
	}

	cpf := baseStr + strconv.Itoa(secondDigit)

	firstChar = cpf[0]
	allSame = true
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != firstChar {
			allSame = false
			break
		}
	}

	if allSame {
		return GenerateValidCPF(seed + 1)
	}

	return cpf
}

func GenerateValidCNPJ(seed int64) string {
	if seed < 0 {
		seed = -seed
	}

	n := seed % 1000000000000
	base := make([]int, 12)
	for i := 11; i >= 0; i-- {
		base[i] = int(n % 10)
		n /= 10
	}

	allSame := true
	for i := 1; i < 12; i++ {
		if base[i] != base[0] {
			allSame = false
			break
		}
	}
	if allSame {
		base[11] = (base[11] + 1) % 10
	}

	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 12; i++ {
		sum += base[i] * weights1[i]
	}
	r := sum % 11
	d1 := 0
	if r >= 2 {
		d1 = 11 - r
	}

	cnpj13 := append(base, d1)
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum = 0
	for i := 0; i < 13; i++ {
		sum += cnpj13[i] * weights2[i]
	}
	r = sum % 11
	d2 := 0
	if r >= 2 {
		d2 = 11 - r
	}

	cnpj14 := append(cnpj13, d2)
	result := ""
	for _, d := range cnpj14 {
		result += strconv.Itoa(d)
	}

	allSame = true
	for i := 1; i < 14; i++ {
		if result[i] != result[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return GenerateValidCNPJ(seed + 1)
	}

	return result
}

func GenerateValidPlate(seed int64) string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	letter1 := letters[seed%26]
	letter2 := letters[(seed/26)%26]
	letter3 := letters[(seed/676)%26]

	numbers := seed % 10000
	if numbers < 0 {
		numbers = -numbers
	}

	numbersStr := strconv.FormatInt(numbers, 10)
	for len(numbersStr) < 4 {
		numbersStr = "0" + numbersStr
	}

	return string(letter1) + string(letter2) + string(letter3) + numbersStr
}
