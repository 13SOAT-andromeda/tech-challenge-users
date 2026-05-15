package services

import (
	"errors"
	"testing"

	"tech-challenge-users/internal/domain"
)

// --- mock company repo ---

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
			delete(m.byDocument, existing.Document)
			clone := *c
			m.companies[i] = &clone
			m.byDocument[c.Document] = &clone
			return nil
		}
	}
	return nil
}

func (m *mockCompanyRepo) Delete(id int64) error {
	for i, c := range m.companies {
		if c.ID == id {
			delete(m.byDocument, c.Document)
			m.companies = append(m.companies[:i], m.companies[i+1:]...)
			return nil
		}
	}
	return nil
}

// --- helpers ---

func buildCompanyService() (*CompanyService, *mockCompanyRepo) {
	repo := newMockCompanyRepo()
	return NewCompanyService(repo), repo
}

// valid CNPJ for tests
const testCNPJ = "11222333000181"

func createTestCompany(svc *CompanyService) (*domain.Company, error) {
	return svc.CreateCompany(CreateCompanyInput{
		Name:     "Empresa X",
		Email:    "contato@empresax.com",
		Document: testCNPJ,
		Contact:  "11999999999",
	})
}

// --- tests ---

func TestCreateCompany_Success(t *testing.T) {
	svc, _ := buildCompanyService()

	c, err := createTestCompany(svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ID == 0 {
		t.Error("expected company ID to be set")
	}
	if c.Document != testCNPJ {
		t.Errorf("expected document %s, got %s", testCNPJ, c.Document)
	}
}

func TestCreateCompany_DuplicateDocument(t *testing.T) {
	svc, _ := buildCompanyService()

	if _, err := createTestCompany(svc); err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	_, err := createTestCompany(svc)
	if !errors.Is(err, ErrDocumentAlreadyExists) {
		t.Errorf("expected ErrDocumentAlreadyExists, got %v", err)
	}
}

func TestCreateCompany_InvalidDocument(t *testing.T) {
	svc, _ := buildCompanyService()

	_, err := svc.CreateCompany(CreateCompanyInput{
		Name:     "Empresa Y",
		Email:    "y@y.com",
		Document: "00000000000000",
		Contact:  "11999999999",
	})
	if err == nil {
		t.Error("expected error for invalid document")
	}
}

func TestGetCompanyByID_NotFound(t *testing.T) {
	svc, _ := buildCompanyService()

	_, err := svc.GetByID(99)
	if !errors.Is(err, ErrCompanyNotFound) {
		t.Errorf("expected ErrCompanyNotFound, got %v", err)
	}
}

func TestUpdateCompany_Success(t *testing.T) {
	svc, _ := buildCompanyService()

	c, _ := createTestCompany(svc)
	newName := "Empresa X Atualizada"

	updated, err := svc.UpdateCompany(c.ID, UpdateCompanyInput{Name: &newName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != newName {
		t.Errorf("expected name %q, got %q", newName, updated.Name)
	}
}

func TestUpdateCompany_NotFound(t *testing.T) {
	svc, _ := buildCompanyService()

	name := "X"
	_, err := svc.UpdateCompany(99, UpdateCompanyInput{Name: &name})
	if !errors.Is(err, ErrCompanyNotFound) {
		t.Errorf("expected ErrCompanyNotFound, got %v", err)
	}
}

func TestDeleteCompany_Success(t *testing.T) {
	svc, repo := buildCompanyService()

	c, _ := createTestCompany(svc)
	if err := svc.DeleteCompany(c.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.companies) != 0 {
		t.Error("expected company to be deleted")
	}
}

func TestDeleteCompany_NotFound(t *testing.T) {
	svc, _ := buildCompanyService()

	err := svc.DeleteCompany(99)
	if !errors.Is(err, ErrCompanyNotFound) {
		t.Errorf("expected ErrCompanyNotFound, got %v", err)
	}
}
