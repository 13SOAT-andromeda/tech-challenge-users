package services

import (
	"testing"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
)

// fullMockUserRepo adds a proper FindAll to the base mock (which always returns nil).
type fullMockUserRepo struct {
	mockUserRepo
}

func (m *fullMockUserRepo) FindAll(_ ports.UserFilters) ([]domain.User, error) {
	out := make([]domain.User, len(m.users))
	for i, u := range m.users {
		out[i] = *u
	}
	return out, nil
}

func buildServiceFull() (*UserService, *mockPersonRepo, *fullMockUserRepo, *mockEmployeeRepo) {
	pr := newMockPersonRepo()
	ur := &fullMockUserRepo{*newMockUserRepo()}
	er := newMockEmployeeRepo()
	tx := &mockTransactor{personRepo: pr, userRepo: ur, employeeRepo: er}
	svc := NewUserService(tx, pr, ur, er)
	return svc, pr, ur, er
}

func TestListUsers_ReturnsAll(t *testing.T) {
	svc, _, ur, _ := buildServiceFull()

	svc.CreateUser(CreateUserInput{
		Name: "Alice", Email: "alice@example.com", Password: "Alice@123",
		Role: domain.RoleAttendant, Document: "52998224725", Position: "Att",
	})
	svc.CreateUser(CreateUserInput{
		Name: "Bob", Email: "bob@example.com", Password: "Bob@12345",
		Role: domain.RoleCustomer, Document: "11144477735",
	})

	outputs, err := svc.ListUsers(ports.UserFilters{})
	if err != nil {
		t.Fatalf("unexpected error listing users: %v", err)
	}
	if len(outputs) != len(ur.users) {
		t.Errorf("expected %d users, got %d", len(ur.users), len(outputs))
	}
}

func TestListUsers_IncludesEmployee(t *testing.T) {
	svc, _, _, _ := buildServiceFull()

	svc.CreateUser(CreateUserInput{
		Name: "Alice", Email: "alice@example.com", Password: "Alice@123",
		Role: domain.RoleAttendant, Document: "52998224725", Position: "Mechanic",
	})

	outputs, err := svc.ListUsers(ports.UserFilters{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(outputs) == 1 && outputs[0].Employee == nil {
		t.Error("expected employee to be set for attendant in list")
	}
}

func TestListUsers_Empty(t *testing.T) {
	svc, _, _, _ := buildService()

	outputs, err := svc.ListUsers(ports.UserFilters{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(outputs) != 0 {
		t.Errorf("expected 0 users, got %d", len(outputs))
	}
}

func TestUpdateUser_Success(t *testing.T) {
	svc, _, _, _ := buildService()

	out, _ := svc.CreateUser(CreateUserInput{
		Name: "Alice", Email: "alice@example.com", Password: "Alice@123",
		Role: domain.RoleAttendant, Document: "52998224725", Position: "Att",
	})

	newName := "Alice Updated"
	newContact := "11999999999"
	updated, err := svc.UpdateUser(out.User.ID, UpdateUserInput{
		Name:    &newName,
		Contact: &newContact,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Person.Name != "Alice Updated" {
		t.Errorf("expected updated name, got %q", updated.Person.Name)
	}
}

func TestUpdateUser_WithEmail(t *testing.T) {
	svc, _, _, _ := buildService()

	out, _ := svc.CreateUser(CreateUserInput{
		Name: "Alice", Email: "alice@example.com", Password: "Alice@123",
		Role: domain.RoleAttendant, Document: "52998224725", Position: "Att",
	})

	newEmail := "alice_new@example.com"
	updated, err := svc.UpdateUser(out.User.ID, UpdateUserInput{Email: &newEmail})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Person.Email != "alice_new@example.com" {
		t.Errorf("expected updated email, got %q", updated.Person.Email)
	}
}

func TestUpdateUser_WithAddress(t *testing.T) {
	svc, _, _, _ := buildService()

	out, _ := svc.CreateUser(CreateUserInput{
		Name: "Alice", Email: "alice@example.com", Password: "Alice@123",
		Role: domain.RoleAttendant, Document: "52998224725",
	})

	addr := domain.Address{Street: "Rua A", Number: "10", City: "SP"}
	_, err := svc.UpdateUser(out.User.ID, UpdateUserInput{Address: &addr})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateUser_DuplicateEmail(t *testing.T) {
	svc, _, _, _ := buildService()

	svc.CreateUser(CreateUserInput{
		Name: "Alice", Email: "alice@example.com", Password: "Alice@123",
		Role: domain.RoleAttendant, Document: "52998224725", Position: "Att",
	})
	out2, _ := svc.CreateUser(CreateUserInput{
		Name: "Bob", Email: "bob@example.com", Password: "Bob@12345",
		Role: domain.RoleCustomer, Document: "11144477735",
	})

	aliceEmail := "alice@example.com"
	_, err := svc.UpdateUser(out2.User.ID, UpdateUserInput{Email: &aliceEmail})
	if err == nil {
		t.Error("expected conflict error for duplicate email")
	}
}

func TestGetByID_Success(t *testing.T) {
	svc, _, _, _ := buildService()

	out, _ := svc.CreateUser(CreateUserInput{
		Name: "Alice", Email: "alice@example.com", Password: "Alice@123",
		Role: domain.RoleAttendant, Document: "52998224725", Position: "Att",
	})

	found, err := svc.GetByID(out.User.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Person.Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %q", found.Person.Email)
	}
	if found.Employee == nil {
		t.Error("expected employee to be set for attendant")
	}
}

func TestGetByDocument_Success(t *testing.T) {
	svc, _, _, _ := buildService()

	svc.CreateUser(CreateUserInput{
		Name: "Alice", Email: "alice@example.com", Password: "Alice@123",
		Role: domain.RoleAttendant, Document: "52998224725",
	})

	out, err := svc.GetByDocument("52998224725")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Person.Document != "52998224725" {
		t.Errorf("expected document '52998224725', got %q", out.Person.Document)
	}
}

func TestGetByDocument_PersonNotFound(t *testing.T) {
	svc, _, _, _ := buildService()

	_, err := svc.GetByDocument("52998224725")
	if err == nil {
		t.Error("expected error when person not found")
	}
}

func TestCreateAdminUser_CreatesOnce(t *testing.T) {
	svc, _, ur, _ := buildService()

	svc.CreateAdminUser(AdminConfig{
		Email:    "admin@example.com",
		Password: "Admin@123",
		Document: "52998224725",
	})

	if len(ur.users) != 1 {
		t.Errorf("expected 1 user after CreateAdminUser, got %d", len(ur.users))
	}

	// calling again should be idempotent
	svc.CreateAdminUser(AdminConfig{
		Email:    "admin@example.com",
		Password: "Admin@123",
		Document: "52998224725",
	})

	if len(ur.users) != 1 {
		t.Errorf("expected still 1 user after duplicate CreateAdminUser, got %d", len(ur.users))
	}
}

func TestDeleteUser_Success(t *testing.T) {
	svc, _, _, _ := buildService()

	out, _ := svc.CreateUser(CreateUserInput{
		Name: "Alice", Email: "alice@example.com", Password: "Alice@123",
		Role: domain.RoleCustomer, Document: "52998224725",
	})

	if err := svc.DeleteUser(out.User.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
