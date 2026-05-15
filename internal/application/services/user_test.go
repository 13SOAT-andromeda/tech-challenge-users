package services

import (
	"errors"
	"testing"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
)

// --- mocks ---

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
func (m *mockPersonRepo) Update(p *domain.Person) error { return nil }

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
func (m *mockUserRepo) GetByEmail(email string) (*domain.User, error)       { return nil, nil }
func (m *mockUserRepo) GetByPersonID(pid int64) (*domain.User, error) {
	for _, u := range m.users {
		if u.PersonID == pid {
			return u, nil
		}
	}
	return nil, nil
}
func (m *mockUserRepo) FindAll(filters ports.UserFilters) ([]domain.User, error) { return nil, nil }
func (m *mockUserRepo) Update(u *domain.User) error                              { return nil }
func (m *mockUserRepo) Delete(id int64) error                                    { return nil }

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

type mockTransactor struct {
	personRepo   ports.PersonRepository
	userRepo     ports.UserRepository
	employeeRepo ports.EmployeeRepository
	customerRepo ports.CustomerRepository
}

func (t *mockTransactor) WithinTransaction(fn func(repos ports.Repositories) error) error {
	return fn(ports.Repositories{
		Person:   t.personRepo,
		User:     t.userRepo,
		Employee: t.employeeRepo,
		Customer: t.customerRepo,
	})
}

// --- helpers ---

func buildService() (*UserService, *mockPersonRepo, *mockUserRepo, *mockEmployeeRepo) {
	pr := newMockPersonRepo()
	ur := newMockUserRepo()
	er := newMockEmployeeRepo()
	tx := &mockTransactor{personRepo: pr, userRepo: ur, employeeRepo: er}
	svc := NewUserService(tx, pr, ur, er)
	return svc, pr, ur, er
}

// --- tests ---

func TestCreateUser_Success(t *testing.T) {
	svc, _, ur, er := buildService()

	out, err := svc.CreateUser(CreateUserInput{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "Alice@123",
		Role:     domain.RoleAttendant,
		Document: "52998224725", // valid CPF
		Position: "Attendant",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.User.ID == 0 {
		t.Error("expected user ID to be set")
	}
	if out.Person.ID == 0 {
		t.Error("expected person ID to be set")
	}
	if out.Employee == nil {
		t.Error("expected employee to be created for attendant role")
	}
	if len(ur.users) != 1 {
		t.Errorf("expected 1 user in repo, got %d", len(ur.users))
	}
	if len(er.employees) != 1 {
		t.Errorf("expected 1 employee in repo, got %d", len(er.employees))
	}
}

func TestCreateUser_CustomerHasNoEmployee(t *testing.T) {
	svc, _, _, er := buildService()

	_, err := svc.CreateUser(CreateUserInput{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "Bob@12345",
		Role:     domain.RoleCustomer,
		Document: "52998224725",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(er.employees) != 0 {
		t.Error("expected no employee for customer role")
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	svc, _, _, _ := buildService()

	input := CreateUserInput{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "Alice@123",
		Role:     domain.RoleAttendant,
		Document: "52998224725",
		Position: "Attendant",
	}
	if _, err := svc.CreateUser(input); err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	input.Document = "11144477735" // different valid CPF
	_, err := svc.CreateUser(input)
	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Errorf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestCreateUser_DuplicateDocument(t *testing.T) {
	svc, _, _, _ := buildService()

	input := CreateUserInput{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "Alice@123",
		Role:     domain.RoleAttendant,
		Document: "52998224725",
		Position: "Attendant",
	}
	if _, err := svc.CreateUser(input); err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	input.Email = "other@example.com"
	_, err := svc.CreateUser(input)
	if !errors.Is(err, ErrDocumentAlreadyExists) {
		t.Errorf("expected ErrDocumentAlreadyExists, got %v", err)
	}
}

func TestCreateUser_InvalidPassword(t *testing.T) {
	svc, _, _, _ := buildService()

	_, err := svc.CreateUser(CreateUserInput{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "weak",
		Role:     domain.RoleAttendant,
		Document: "52998224725",
	})
	if err == nil {
		t.Error("expected validation error for weak password")
	}
}

func TestCreateUser_InvalidRole(t *testing.T) {
	svc, _, _, _ := buildService()

	_, err := svc.CreateUser(CreateUserInput{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "Alice@123",
		Role:     "superadmin",
		Document: "52998224725",
	})
	if err == nil {
		t.Error("expected error for invalid role")
	}
}

func TestGetByID_NotFound(t *testing.T) {
	svc, _, _, _ := buildService()

	_, err := svc.GetByID(99)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	svc, _, _, _ := buildService()

	err := svc.DeleteUser(99)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
