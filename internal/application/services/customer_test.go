package services

import (
	"errors"
	"testing"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
)

// --- mock customer repo ---

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
func (m *mockCustomerRepo) Update(c *domain.Customer) error { return nil }
func (m *mockCustomerRepo) Delete(_ int64) error            { return nil }

// failingUserRepo always returns an error on Create (rollback test)
type failingUserRepo struct{}

func (r *failingUserRepo) Create(_ *domain.User) error                          { return errors.New("db error") }
func (r *failingUserRepo) FindByID(_ int64) (*domain.User, error)               { return nil, nil }
func (r *failingUserRepo) GetByEmail(_ string) (*domain.User, error)            { return nil, nil }
func (r *failingUserRepo) GetByPersonID(_ int64) (*domain.User, error)          { return nil, nil }
func (r *failingUserRepo) FindAll(_ ports.UserFilters) ([]domain.User, error)   { return nil, nil }
func (r *failingUserRepo) Update(_ *domain.User) error                          { return nil }
func (r *failingUserRepo) Delete(_ int64) error                                 { return nil }

// --- helpers ---

func buildCustomerService() (*CustomerService, *mockPersonRepo, *mockCustomerRepo) {
	pr := newMockPersonRepo()
	ur := newMockUserRepo()
	er := newMockEmployeeRepo()
	cr := newMockCustomerRepo()
	tx := &mockTransactor{personRepo: pr, userRepo: ur, employeeRepo: er, customerRepo: cr}
	svc := NewCustomerService(tx, pr, cr)
	return svc, pr, cr
}

// --- tests ---

func TestCreateCustomer_Success(t *testing.T) {
	svc, _, cr := buildCustomerService()

	out, err := svc.CreateCustomer(CreateCustomerInput{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "Bob@12345",
		Type:     domain.CustomerTypePF,
		Document: "52998224725",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Customer.ID == 0 {
		t.Error("expected customer ID to be set")
	}
	if out.Customer.Type != domain.CustomerTypePF {
		t.Errorf("expected type pf, got %s", out.Customer.Type)
	}
	if len(cr.customers) != 1 {
		t.Errorf("expected 1 customer in repo, got %d", len(cr.customers))
	}
}

func TestCreateCustomer_InvalidType(t *testing.T) {
	svc, _, _ := buildCustomerService()

	_, err := svc.CreateCustomer(CreateCustomerInput{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "Bob@12345",
		Type:     "corporate",
		Document: "52998224725",
	})
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestCreateCustomer_DuplicateEmail(t *testing.T) {
	svc, _, _ := buildCustomerService()

	input := CreateCustomerInput{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "Bob@12345",
		Type:     domain.CustomerTypePF,
		Document: "52998224725",
	}
	if _, err := svc.CreateCustomer(input); err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	input.Document = "11144477735"
	_, err := svc.CreateCustomer(input)
	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Errorf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestCreateCustomer_RollbackOnUserCreateError(t *testing.T) {
	pr := newMockPersonRepo()
	cr := newMockCustomerRepo()
	tx := &mockTransactor{
		personRepo:   pr,
		userRepo:     &failingUserRepo{},
		employeeRepo: newMockEmployeeRepo(),
		customerRepo: cr,
	}
	svc := NewCustomerService(tx, pr, cr)

	_, err := svc.CreateCustomer(CreateCustomerInput{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "Bob@12345",
		Type:     domain.CustomerTypePF,
		Document: "52998224725",
	})
	if err == nil {
		t.Error("expected error when user creation fails")
	}
	if len(cr.customers) != 0 {
		t.Error("customer should not be created when transaction fails")
	}
}

func TestGetCustomerByID_NotFound(t *testing.T) {
	svc, _, _ := buildCustomerService()

	_, err := svc.GetByID(99)
	if !errors.Is(err, ErrCustomerNotFound) {
		t.Errorf("expected ErrCustomerNotFound, got %v", err)
	}
}

func TestDeleteCustomer_NotFound(t *testing.T) {
	svc, _, _ := buildCustomerService()

	err := svc.DeleteCustomer(99)
	if !errors.Is(err, ErrCustomerNotFound) {
		t.Errorf("expected ErrCustomerNotFound, got %v", err)
	}
}
