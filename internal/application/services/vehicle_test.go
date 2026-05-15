package services

import (
	"errors"
	"testing"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
)

// --- mock vehicle repo ---

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
			// update byPlate index
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

// --- helpers ---

func buildVehicleService() (*VehicleService, *mockVehicleRepo) {
	repo := newMockVehicleRepo()
	return NewVehicleService(repo), repo
}

func createTestVehicle(svc *VehicleService, plate string) (*domain.Vehicle, error) {
	return svc.CreateVehicle(CreateVehicleInput{
		Plate: plate,
		Name:  "Civic",
		Year:  2022,
		Brand: "Honda",
		Color: "Prata",
	})
}

// --- tests ---

func TestCreateVehicle_Success(t *testing.T) {
	svc, repo := buildVehicleService()

	v, err := createTestVehicle(svc, "ABC1234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.ID == 0 {
		t.Error("expected vehicle ID to be set")
	}
	if v.Plate != "ABC1234" {
		t.Errorf("expected plate ABC1234, got %s", v.Plate)
	}
	if len(repo.vehicles) != 1 {
		t.Errorf("expected 1 vehicle in repo, got %d", len(repo.vehicles))
	}
}

func TestCreateVehicle_MercosulPlate(t *testing.T) {
	svc, _ := buildVehicleService()

	v, err := createTestVehicle(svc, "ABC1D23")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Plate != "ABC1D23" {
		t.Errorf("expected plate ABC1D23, got %s", v.Plate)
	}
}

func TestCreateVehicle_InvalidPlate(t *testing.T) {
	svc, _ := buildVehicleService()

	_, err := createTestVehicle(svc, "INVALID")
	if err == nil {
		t.Error("expected error for invalid plate")
	}
}

func TestCreateVehicle_DuplicatePlate(t *testing.T) {
	svc, _ := buildVehicleService()

	if _, err := createTestVehicle(svc, "ABC1234"); err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	_, err := createTestVehicle(svc, "ABC1234")
	if !errors.Is(err, ErrPlateAlreadyExists) {
		t.Errorf("expected ErrPlateAlreadyExists, got %v", err)
	}
}

func TestUpdateVehicle_SamePlateAllowed(t *testing.T) {
	svc, _ := buildVehicleService()

	v, _ := createTestVehicle(svc, "ABC1234")

	updated, err := svc.UpdateVehicle(v.ID, UpdateVehicleInput{
		Plate: strPtr("ABC1234"),
		Color: strPtr("Preto"),
	})
	if err != nil {
		t.Fatalf("expected no error updating with same plate, got %v", err)
	}
	if updated.Color != "Preto" {
		t.Errorf("expected color Preto, got %s", updated.Color)
	}
}

func TestUpdateVehicle_DuplicatePlateOtherVehicle(t *testing.T) {
	svc, _ := buildVehicleService()

	createTestVehicle(svc, "ABC1234")
	v2, _ := createTestVehicle(svc, "XYZ9876")

	_, err := svc.UpdateVehicle(v2.ID, UpdateVehicleInput{Plate: strPtr("ABC1234")})
	if !errors.Is(err, ErrPlateAlreadyExists) {
		t.Errorf("expected ErrPlateAlreadyExists, got %v", err)
	}
}

func TestGetVehicleByID_NotFound(t *testing.T) {
	svc, _ := buildVehicleService()

	_, err := svc.GetByID(99)
	if !errors.Is(err, ErrVehicleNotFound) {
		t.Errorf("expected ErrVehicleNotFound, got %v", err)
	}
}

func TestDeleteVehicle_NotFound(t *testing.T) {
	svc, _ := buildVehicleService()

	err := svc.DeleteVehicle(99)
	if !errors.Is(err, ErrVehicleNotFound) {
		t.Errorf("expected ErrVehicleNotFound, got %v", err)
	}
}

func strPtr(s string) *string { return &s }
