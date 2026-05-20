package services

import (
	"testing"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
)

func TestListCustomers_ReturnsAll(t *testing.T) {
	svc, _, _ := buildCustomerService()

	svc.CreateCustomer(CreateCustomerInput{
		Name: "Bob", Email: "bob@example.com", Password: "Bob@12345",
		Type: domain.CustomerTypePF, Document: "52998224725",
	})
	svc.CreateCustomer(CreateCustomerInput{
		Name: "Carol", Email: "carol@example.com", Password: "Carol@123",
		Type: domain.CustomerTypePJ, Document: "11144477735",
	})

	outputs, err := svc.ListCustomers(ports.CustomerFilters{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(outputs) != 2 {
		t.Errorf("expected 2 customers, got %d", len(outputs))
	}
}

func TestListCustomers_Empty(t *testing.T) {
	svc, _, _ := buildCustomerService()

	outputs, err := svc.ListCustomers(ports.CustomerFilters{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(outputs) != 0 {
		t.Errorf("expected 0 customers, got %d", len(outputs))
	}
}

func TestGetCustomerByID_Success(t *testing.T) {
	svc, _, _ := buildCustomerService()

	out, _ := svc.CreateCustomer(CreateCustomerInput{
		Name: "Bob", Email: "bob@example.com", Password: "Bob@12345",
		Type: domain.CustomerTypePF, Document: "52998224725",
	})

	found, err := svc.GetByID(out.Customer.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Person.Email != "bob@example.com" {
		t.Errorf("expected email bob@example.com, got %q", found.Person.Email)
	}
}

func TestUpdateCustomer_Name(t *testing.T) {
	svc, _, _ := buildCustomerService()

	out, _ := svc.CreateCustomer(CreateCustomerInput{
		Name: "Bob", Email: "bob@example.com", Password: "Bob@12345",
		Type: domain.CustomerTypePF, Document: "52998224725",
	})

	newName := "Bob Updated"
	updated, err := svc.UpdateCustomer(out.Customer.ID, UpdateCustomerInput{Name: &newName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Person.Name != "Bob Updated" {
		t.Errorf("expected 'Bob Updated', got %q", updated.Person.Name)
	}
}

func TestUpdateCustomer_Email(t *testing.T) {
	svc, _, _ := buildCustomerService()

	out, _ := svc.CreateCustomer(CreateCustomerInput{
		Name: "Bob", Email: "bob@example.com", Password: "Bob@12345",
		Type: domain.CustomerTypePF, Document: "52998224725",
	})

	newEmail := "bob_new@example.com"
	updated, err := svc.UpdateCustomer(out.Customer.ID, UpdateCustomerInput{Email: &newEmail})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Person.Email != "bob_new@example.com" {
		t.Errorf("expected updated email, got %q", updated.Person.Email)
	}
}

func TestUpdateCustomer_DuplicateEmail(t *testing.T) {
	svc, _, _ := buildCustomerService()

	svc.CreateCustomer(CreateCustomerInput{
		Name: "Bob", Email: "bob@example.com", Password: "Bob@12345",
		Type: domain.CustomerTypePF, Document: "52998224725",
	})
	out2, _ := svc.CreateCustomer(CreateCustomerInput{
		Name: "Carol", Email: "carol@example.com", Password: "Carol@123",
		Type: domain.CustomerTypePF, Document: "11144477735",
	})

	bobEmail := "bob@example.com"
	_, err := svc.UpdateCustomer(out2.Customer.ID, UpdateCustomerInput{Email: &bobEmail})
	if err == nil {
		t.Error("expected conflict for duplicate email")
	}
}

func TestUpdateCustomer_Type(t *testing.T) {
	svc, _, _ := buildCustomerService()

	out, _ := svc.CreateCustomer(CreateCustomerInput{
		Name: "Bob", Email: "bob@example.com", Password: "Bob@12345",
		Type: domain.CustomerTypePF, Document: "52998224725",
	})

	newType := domain.CustomerTypePJ
	updated, err := svc.UpdateCustomer(out.Customer.ID, UpdateCustomerInput{Type: &newType})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Customer.Type != domain.CustomerTypePJ {
		t.Errorf("expected type pj, got %q", updated.Customer.Type)
	}
}

func TestUpdateCustomer_InvalidType(t *testing.T) {
	svc, _, _ := buildCustomerService()

	out, _ := svc.CreateCustomer(CreateCustomerInput{
		Name: "Bob", Email: "bob@example.com", Password: "Bob@12345",
		Type: domain.CustomerTypePF, Document: "52998224725",
	})

	badType := "invalid"
	_, err := svc.UpdateCustomer(out.Customer.ID, UpdateCustomerInput{Type: &badType})
	if err == nil {
		t.Error("expected error for invalid customer type")
	}
}

func TestUpdateCustomer_WithAddress(t *testing.T) {
	svc, _, _ := buildCustomerService()

	out, _ := svc.CreateCustomer(CreateCustomerInput{
		Name: "Bob", Email: "bob@example.com", Password: "Bob@12345",
		Type: domain.CustomerTypePF, Document: "52998224725",
	})

	addr := domain.Address{Street: "Rua B", Number: "5", City: "RJ"}
	_, err := svc.UpdateCustomer(out.Customer.ID, UpdateCustomerInput{Address: &addr})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateCustomer_NotFound(t *testing.T) {
	svc, _, _ := buildCustomerService()

	name := "X"
	_, err := svc.UpdateCustomer(999, UpdateCustomerInput{Name: &name})
	if err == nil {
		t.Error("expected error for non-existent customer")
	}
}

func TestDeleteCustomer_Success(t *testing.T) {
	svc, _, _ := buildCustomerService()

	out, _ := svc.CreateCustomer(CreateCustomerInput{
		Name: "Bob", Email: "bob@example.com", Password: "Bob@12345",
		Type: domain.CustomerTypePF, Document: "52998224725",
	})

	if err := svc.DeleteCustomer(out.Customer.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
