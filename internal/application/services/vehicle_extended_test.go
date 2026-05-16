package services

import (
	"testing"

	"tech-challenge-users/internal/application/ports"
)

func TestListVehicles_ReturnsAll(t *testing.T) {
	svc, _ := buildVehicleService()

	createTestVehicle(svc, "ABC1234")
	createTestVehicle(svc, "XYZ9876")

	vehicles, err := svc.ListVehicles(ports.VehicleFilters{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vehicles) != 2 {
		t.Errorf("expected 2 vehicles, got %d", len(vehicles))
	}
}

func TestListVehicles_Empty(t *testing.T) {
	svc, _ := buildVehicleService()

	vehicles, err := svc.ListVehicles(ports.VehicleFilters{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vehicles) != 0 {
		t.Errorf("expected 0 vehicles, got %d", len(vehicles))
	}
}

func TestGetVehicleByID_Success(t *testing.T) {
	svc, _ := buildVehicleService()

	v, _ := createTestVehicle(svc, "ABC1234")

	found, err := svc.GetByID(v.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Plate != "ABC1234" {
		t.Errorf("expected plate 'ABC1234', got %q", found.Plate)
	}
}

func TestUpdateVehicle_AllFields(t *testing.T) {
	svc, _ := buildVehicleService()

	v, _ := createTestVehicle(svc, "ABC1234")

	newYear := 2023
	updated, err := svc.UpdateVehicle(v.ID, UpdateVehicleInput{
		Name:  strPtr("Fit"),
		Year:  &newYear,
		Brand: strPtr("Honda"),
		Color: strPtr("Azul"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Fit" || updated.Year != 2023 || updated.Color != "Azul" {
		t.Error("expected updated fields")
	}
}

func TestUpdateVehicle_NotFound(t *testing.T) {
	svc, _ := buildVehicleService()

	_, err := svc.UpdateVehicle(999, UpdateVehicleInput{Color: strPtr("Preto")})
	if err == nil {
		t.Error("expected error for non-existent vehicle")
	}
}

func TestDeleteVehicle_Success(t *testing.T) {
	svc, repo := buildVehicleService()

	v, _ := createTestVehicle(svc, "ABC1234")
	if err := svc.DeleteVehicle(v.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.vehicles) != 0 {
		t.Error("expected vehicle to be deleted")
	}
}
