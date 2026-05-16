package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"

	"github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

// ---------------------------------------------------------------------------
// Shared test helpers
// ---------------------------------------------------------------------------

func setupRouter(fn func(*gin.RouterGroup)) *gin.Engine {
	r := gin.New()
	fn(r.Group("/"))
	return r
}

func doRequest(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	r.ServeHTTP(w, req)
	return w
}

func sp(s string) *string { return &s }
func ip(i int) *int       { return &i }

// ---------------------------------------------------------------------------
// mockPersonRepo
// ---------------------------------------------------------------------------

type mockPersonRepo struct {
	byEmail    map[string]*domain.Person
	byDocument map[string]*domain.Person
	created    []*domain.Person
	nextID     int64
}

func newMockPersonRepo() *mockPersonRepo {
	return &mockPersonRepo{
		byEmail:    make(map[string]*domain.Person),
		byDocument: make(map[string]*domain.Person),
		nextID:     1,
	}
}

func (m *mockPersonRepo) Create(p *domain.Person) error {
	p.ID = m.nextID
	m.nextID++
	clone := *p
	m.created = append(m.created, &clone)
	m.byEmail[p.Email] = &clone
	m.byDocument[p.Document] = &clone
	return nil
}
func (m *mockPersonRepo) FindByID(id int64) (*domain.Person, error) {
	for _, p := range m.created {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, nil
}
func (m *mockPersonRepo) GetByEmail(email string) (*domain.Person, error) {
	return m.byEmail[email], nil
}
func (m *mockPersonRepo) GetByDocument(doc string) (*domain.Person, error) {
	return m.byDocument[doc], nil
}
func (m *mockPersonRepo) Update(p *domain.Person) error {
	for i, existing := range m.created {
		if existing.ID == p.ID {
			clone := *p
			m.created[i] = &clone
			m.byEmail[p.Email] = &clone
			return nil
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// mockUserRepo
// ---------------------------------------------------------------------------

type mockUserRepo struct {
	users  []*domain.User
	nextID int64
}

func newMockUserRepo() *mockUserRepo { return &mockUserRepo{nextID: 1} }

func (m *mockUserRepo) Create(u *domain.User) error {
	u.ID = m.nextID
	m.nextID++
	clone := *u
	m.users = append(m.users, &clone)
	return nil
}
func (m *mockUserRepo) FindByID(id int64) (*domain.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}
func (m *mockUserRepo) GetByEmail(_ string) (*domain.User, error) { return nil, nil }
func (m *mockUserRepo) GetByPersonID(pid int64) (*domain.User, error) {
	for _, u := range m.users {
		if u.PersonID == pid {
			return u, nil
		}
	}
	return nil, nil
}
func (m *mockUserRepo) FindAll(_ ports.UserFilters) ([]domain.User, error) {
	out := make([]domain.User, len(m.users))
	for i, u := range m.users {
		out[i] = *u
	}
	return out, nil
}
func (m *mockUserRepo) Update(_ *domain.User) error { return nil }
func (m *mockUserRepo) Delete(id int64) error {
	for i, u := range m.users {
		if u.ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// mockEmployeeRepo
// ---------------------------------------------------------------------------

type mockEmployeeRepo struct {
	employees []*domain.Employee
	nextID    int64
}

func newMockEmployeeRepo() *mockEmployeeRepo { return &mockEmployeeRepo{nextID: 1} }

func (m *mockEmployeeRepo) Create(e *domain.Employee) error {
	e.ID = m.nextID
	m.nextID++
	clone := *e
	m.employees = append(m.employees, &clone)
	return nil
}
func (m *mockEmployeeRepo) FindByPersonID(pid int64) (*domain.Employee, error) {
	for _, e := range m.employees {
		if e.PersonID == pid {
			return e, nil
		}
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// mockCustomerRepo
// ---------------------------------------------------------------------------

type mockCustomerRepo struct {
	customers []*domain.Customer
	nextID    int64
}

func newMockCustomerRepo() *mockCustomerRepo { return &mockCustomerRepo{nextID: 1} }

func (m *mockCustomerRepo) Create(c *domain.Customer) error {
	c.ID = m.nextID
	m.nextID++
	clone := *c
	m.customers = append(m.customers, &clone)
	return nil
}
func (m *mockCustomerRepo) FindByID(id int64) (*domain.Customer, error) {
	for _, c := range m.customers {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, nil
}
func (m *mockCustomerRepo) FindAll(_ ports.CustomerFilters) ([]domain.Customer, error) {
	out := make([]domain.Customer, len(m.customers))
	for i, c := range m.customers {
		out[i] = *c
	}
	return out, nil
}
func (m *mockCustomerRepo) Update(c *domain.Customer) error {
	for i, existing := range m.customers {
		if existing.ID == c.ID {
			clone := *c
			m.customers[i] = &clone
			return nil
		}
	}
	return nil
}
func (m *mockCustomerRepo) Delete(_ int64) error { return nil }

// ---------------------------------------------------------------------------
// mockVehicleRepo
// ---------------------------------------------------------------------------

type mockVehicleRepo struct {
	vehicles []*domain.Vehicle
	byPlate  map[string]*domain.Vehicle
	nextID   int64
}

func newMockVehicleRepo() *mockVehicleRepo {
	return &mockVehicleRepo{byPlate: make(map[string]*domain.Vehicle), nextID: 1}
}

func (m *mockVehicleRepo) Create(v *domain.Vehicle) error {
	v.ID = m.nextID
	m.nextID++
	clone := *v
	m.vehicles = append(m.vehicles, &clone)
	m.byPlate[v.Plate] = &clone
	return nil
}
func (m *mockVehicleRepo) FindByID(id int64) (*domain.Vehicle, error) {
	for _, v := range m.vehicles {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, nil
}
func (m *mockVehicleRepo) FindByPlate(plate string) (*domain.Vehicle, error) {
	return m.byPlate[plate], nil
}
func (m *mockVehicleRepo) FindAll(_ ports.VehicleFilters) ([]domain.Vehicle, error) {
	out := make([]domain.Vehicle, len(m.vehicles))
	for i, v := range m.vehicles {
		out[i] = *v
	}
	return out, nil
}
func (m *mockVehicleRepo) Update(v *domain.Vehicle) error {
	for i, existing := range m.vehicles {
		if existing.ID == v.ID {
			delete(m.byPlate, existing.Plate)
			clone := *v
			m.vehicles[i] = &clone
			m.byPlate[v.Plate] = &clone
			return nil
		}
	}
	return nil
}
func (m *mockVehicleRepo) Delete(id int64) error {
	for i, v := range m.vehicles {
		if v.ID == id {
			delete(m.byPlate, v.Plate)
			m.vehicles = append(m.vehicles[:i], m.vehicles[i+1:]...)
			return nil
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// mockCompanyRepo
// ---------------------------------------------------------------------------

type mockCompanyRepo struct {
	companies  []*domain.Company
	byDocument map[string]*domain.Company
	nextID     int64
}

func newMockCompanyRepo() *mockCompanyRepo {
	return &mockCompanyRepo{byDocument: make(map[string]*domain.Company), nextID: 1}
}

func (m *mockCompanyRepo) Create(c *domain.Company) error {
	c.ID = m.nextID
	m.nextID++
	clone := *c
	m.companies = append(m.companies, &clone)
	m.byDocument[c.Document] = &clone
	return nil
}
func (m *mockCompanyRepo) FindByID(id int64) (*domain.Company, error) {
	for _, c := range m.companies {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, nil
}
func (m *mockCompanyRepo) FindByDocument(doc string) (*domain.Company, error) {
	return m.byDocument[doc], nil
}
func (m *mockCompanyRepo) Update(c *domain.Company) error {
	for i, existing := range m.companies {
		if existing.ID == c.ID {
			clone := *c
			m.companies[i] = &clone
			return nil
		}
	}
	return nil
}
func (m *mockCompanyRepo) Delete(id int64) error {
	for i, c := range m.companies {
		if c.ID == id {
			m.companies = append(m.companies[:i], m.companies[i+1:]...)
			return nil
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// mockCVRepo (CustomerVehicle)
// ---------------------------------------------------------------------------

type mockCVRepo struct {
	cvs    []*domain.CustomerVehicle
	nextID int64
}

func newMockCVRepo() *mockCVRepo { return &mockCVRepo{nextID: 1} }

func (m *mockCVRepo) Associate(cv *domain.CustomerVehicle) error {
	cv.ID = m.nextID
	m.nextID++
	clone := *cv
	m.cvs = append(m.cvs, &clone)
	return nil
}
func (m *mockCVRepo) Dissociate(customerID, vehicleID int64) error {
	for i, cv := range m.cvs {
		if cv.CustomerID == customerID && cv.VehicleID == vehicleID {
			m.cvs = append(m.cvs[:i], m.cvs[i+1:]...)
			return nil
		}
	}
	return nil
}
func (m *mockCVRepo) FindByCustomerID(customerID int64) ([]domain.CustomerVehicle, error) {
	var out []domain.CustomerVehicle
	for _, cv := range m.cvs {
		if cv.CustomerID == customerID {
			out = append(out, *cv)
		}
	}
	return out, nil
}
func (m *mockCVRepo) FindByCustomerAndVehicle(customerID, vehicleID int64) (*domain.CustomerVehicle, error) {
	for _, cv := range m.cvs {
		if cv.CustomerID == customerID && cv.VehicleID == vehicleID {
			return cv, nil
		}
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// mockTransactor
// ---------------------------------------------------------------------------

type mockTransactor struct {
	pr ports.PersonRepository
	ur ports.UserRepository
	er ports.EmployeeRepository
	cr ports.CustomerRepository
}

func (t *mockTransactor) WithinTransaction(fn func(repos ports.Repositories) error) error {
	return fn(ports.Repositories{
		Person:   t.pr,
		User:     t.ur,
		Employee: t.er,
		Customer: t.cr,
	})
}
