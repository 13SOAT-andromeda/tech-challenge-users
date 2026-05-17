package services

import (
	"testing"

	"tech-challenge-users/internal/domain"
)

func TestEmployeeService_GetByID_Success(t *testing.T) {
	repo := newMockEmployeeRepo()
	_ = repo.Create(&domain.Employee{Position: "mechanic", PersonID: 5})

	svc := NewEmployeeService(repo)
	got, err := svc.GetByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected employee, got nil")
	}
	if got.ID != 1 {
		t.Errorf("expected ID=1, got %d", got.ID)
	}
	if got.Position != "mechanic" {
		t.Errorf("expected position=mechanic, got %s", got.Position)
	}
}

func TestEmployeeService_GetByID_NotFound(t *testing.T) {
	repo := newMockEmployeeRepo()
	svc := NewEmployeeService(repo)

	_, err := svc.GetByID(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != ErrEmployeeNotFound {
		t.Errorf("expected ErrEmployeeNotFound, got %v", err)
	}
}
