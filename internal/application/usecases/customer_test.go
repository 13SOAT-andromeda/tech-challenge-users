package usecases_test

import (
	"errors"
	"testing"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/application/usecases"
	"tech-challenge-users/internal/domain"
)

// --- mock repos ---

type mockCustomerRepo struct {
	findByID func(id int64) (*domain.Customer, error)
}

func (m *mockCustomerRepo) Create(_ *domain.Customer) error { return nil }
func (m *mockCustomerRepo) FindByID(id int64) (*domain.Customer, error) {
	return m.findByID(id)
}
func (m *mockCustomerRepo) FindAll(_ ports.CustomerFilters) ([]domain.Customer, error) {
	return nil, nil
}
func (m *mockCustomerRepo) Update(_ *domain.Customer) error { return nil }
func (m *mockCustomerRepo) Delete(_ int64) error            { return nil }

type mockVehicleRepo struct {
	findByID func(id int64) (*domain.Vehicle, error)
}

func (m *mockVehicleRepo) Create(_ *domain.Vehicle) error { return nil }
func (m *mockVehicleRepo) FindByID(id int64) (*domain.Vehicle, error) {
	return m.findByID(id)
}
func (m *mockVehicleRepo) FindByPlate(_ string) (*domain.Vehicle, error) { return nil, nil }
func (m *mockVehicleRepo) FindAll(_ ports.VehicleFilters) ([]domain.Vehicle, error) {
	return nil, nil
}
func (m *mockVehicleRepo) Update(_ *domain.Vehicle) error { return nil }
func (m *mockVehicleRepo) Delete(_ int64) error            { return nil }

type mockCVRepo struct {
	associate                func(cv *domain.CustomerVehicle) error
	dissociate               func(customerID, vehicleID int64) error
	findByCustomerID         func(customerID int64) ([]domain.CustomerVehicle, error)
	findByCustomerAndVehicle func(customerID, vehicleID int64) (*domain.CustomerVehicle, error)
}

func (m *mockCVRepo) Associate(cv *domain.CustomerVehicle) error { return m.associate(cv) }
func (m *mockCVRepo) Dissociate(cID, vID int64) error            { return m.dissociate(cID, vID) }
func (m *mockCVRepo) FindByCustomerID(cID int64) ([]domain.CustomerVehicle, error) {
	return m.findByCustomerID(cID)
}
func (m *mockCVRepo) FindByCustomerAndVehicle(cID, vID int64) (*domain.CustomerVehicle, error) {
	return m.findByCustomerAndVehicle(cID, vID)
}

// --- helpers ---

func customerFound() func(int64) (*domain.Customer, error) {
	return func(id int64) (*domain.Customer, error) {
		return &domain.Customer{ID: id}, nil
	}
}

func vehicleFound() func(int64) (*domain.Vehicle, error) {
	return func(id int64) (*domain.Vehicle, error) {
		return &domain.Vehicle{ID: id}, nil
	}
}

func noAssociation() func(int64, int64) (*domain.CustomerVehicle, error) {
	return func(_, _ int64) (*domain.CustomerVehicle, error) { return nil, nil }
}

func existingAssociation() func(int64, int64) (*domain.CustomerVehicle, error) {
	return func(cID, vID int64) (*domain.CustomerVehicle, error) {
		return &domain.CustomerVehicle{CustomerID: cID, VehicleID: vID}, nil
	}
}

// --- AddVehicle ---

func TestAddVehicle_Success(t *testing.T) {
	uc := usecases.NewCustomerUseCase(
		&mockCustomerRepo{findByID: customerFound()},
		&mockVehicleRepo{findByID: vehicleFound()},
		&mockCVRepo{
			findByCustomerAndVehicle: noAssociation(),
			associate:                func(_ *domain.CustomerVehicle) error { return nil },
		},
	)

	cv, err := uc.AddVehicle(1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cv.CustomerID != 1 || cv.VehicleID != 2 {
		t.Fatalf("unexpected cv: %+v", cv)
	}
}

func TestAddVehicle_CustomerNotFound(t *testing.T) {
	uc := usecases.NewCustomerUseCase(
		&mockCustomerRepo{findByID: func(_ int64) (*domain.Customer, error) { return nil, nil }},
		&mockVehicleRepo{findByID: vehicleFound()},
		&mockCVRepo{},
	)

	_, err := uc.AddVehicle(1, 2)
	if !errors.Is(err, usecases.ErrCustomerNotFound) {
		t.Fatalf("expected ErrCustomerNotFound, got %v", err)
	}
}

func TestAddVehicle_VehicleNotFound(t *testing.T) {
	uc := usecases.NewCustomerUseCase(
		&mockCustomerRepo{findByID: customerFound()},
		&mockVehicleRepo{findByID: func(_ int64) (*domain.Vehicle, error) { return nil, nil }},
		&mockCVRepo{},
	)

	_, err := uc.AddVehicle(1, 2)
	if !errors.Is(err, usecases.ErrVehicleNotFound) {
		t.Fatalf("expected ErrVehicleNotFound, got %v", err)
	}
}

func TestAddVehicle_AlreadyExists(t *testing.T) {
	uc := usecases.NewCustomerUseCase(
		&mockCustomerRepo{findByID: customerFound()},
		&mockVehicleRepo{findByID: vehicleFound()},
		&mockCVRepo{findByCustomerAndVehicle: existingAssociation()},
	)

	_, err := uc.AddVehicle(1, 2)
	if !errors.Is(err, usecases.ErrAssociationAlreadyExists) {
		t.Fatalf("expected ErrAssociationAlreadyExists, got %v", err)
	}
}

// --- RemoveVehicle ---

func TestRemoveVehicle_Success(t *testing.T) {
	uc := usecases.NewCustomerUseCase(
		&mockCustomerRepo{findByID: customerFound()},
		&mockVehicleRepo{findByID: vehicleFound()},
		&mockCVRepo{
			findByCustomerAndVehicle: existingAssociation(),
			dissociate:               func(_, _ int64) error { return nil },
		},
	)

	if err := uc.RemoveVehicle(1, 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemoveVehicle_NotFound(t *testing.T) {
	uc := usecases.NewCustomerUseCase(
		&mockCustomerRepo{findByID: customerFound()},
		&mockVehicleRepo{findByID: vehicleFound()},
		&mockCVRepo{findByCustomerAndVehicle: noAssociation()},
	)

	err := uc.RemoveVehicle(1, 2)
	if !errors.Is(err, usecases.ErrAssociationNotFound) {
		t.Fatalf("expected ErrAssociationNotFound, got %v", err)
	}
}

// --- GetCustomerVehicles ---

func TestGetCustomerVehicles_Success(t *testing.T) {
	expected := []domain.CustomerVehicle{
		{CustomerID: 1, VehicleID: 10},
		{CustomerID: 1, VehicleID: 11},
	}

	uc := usecases.NewCustomerUseCase(
		&mockCustomerRepo{findByID: customerFound()},
		&mockVehicleRepo{findByID: vehicleFound()},
		&mockCVRepo{
			findByCustomerID: func(_ int64) ([]domain.CustomerVehicle, error) { return expected, nil },
		},
	)

	got, err := uc.GetCustomerVehicles(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 vehicles, got %d", len(got))
	}
}

func TestGetCustomerVehicles_CustomerNotFound(t *testing.T) {
	uc := usecases.NewCustomerUseCase(
		&mockCustomerRepo{findByID: func(_ int64) (*domain.Customer, error) { return nil, nil }},
		&mockVehicleRepo{findByID: vehicleFound()},
		&mockCVRepo{},
	)

	_, err := uc.GetCustomerVehicles(1)
	if !errors.Is(err, usecases.ErrCustomerNotFound) {
		t.Fatalf("expected ErrCustomerNotFound, got %v", err)
	}
}
