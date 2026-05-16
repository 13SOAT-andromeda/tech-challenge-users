package converters

import (
	"testing"
	"time"

	companymodel "tech-challenge-users/internal/adapter/database/model/company"
	customermodel "tech-challenge-users/internal/adapter/database/model/customer"
	cvmodel "tech-challenge-users/internal/adapter/database/model/customer_vehicle"
	employeemodel "tech-challenge-users/internal/adapter/database/model/employee"
	personmodel "tech-challenge-users/internal/adapter/database/model/person"
	usermodel "tech-challenge-users/internal/adapter/database/model/user"
	vehiclemodel "tech-challenge-users/internal/adapter/database/model/vehicle"
	"tech-challenge-users/internal/domain"
)

// --- helpers ---

func strp(s string) *string { return &s }

// --- Person ---

func TestPersonToModel_WithAllFields(t *testing.T) {
	p := domain.Person{
		ID:       1,
		Name:     "Alice",
		Email:    "alice@example.com",
		Document: "52998224725",
		IsActive: true,
		Contact:  "11999999999",
		Address: domain.Address{
			Street:       "Rua A",
			Number:       "10",
			Neighborhood: "Centro",
			City:         "SP",
			Country:      "BR",
			ZipCode:      "01001000",
		},
	}

	m := PersonToModel(p)

	if m.ID != 1 || m.Name != "Alice" || m.Email != "alice@example.com" {
		t.Error("basic fields mismatch")
	}
	if m.Contact == nil || *m.Contact != "11999999999" {
		t.Error("expected Contact to be set")
	}
	if m.Address == nil || *m.Address != "Rua A" {
		t.Error("expected Address to be set")
	}
	if m.AddressNumber == nil || *m.AddressNumber != "10" {
		t.Error("expected AddressNumber to be set")
	}
	if m.Neighborhood == nil || *m.Neighborhood != "Centro" {
		t.Error("expected Neighborhood to be set")
	}
	if m.City == nil || *m.City != "SP" {
		t.Error("expected City to be set")
	}
	if m.Country == nil || *m.Country != "BR" {
		t.Error("expected Country to be set")
	}
	if m.ZipCode == nil || *m.ZipCode != "01001000" {
		t.Error("expected ZipCode to be set")
	}
}

func TestPersonToModel_WithNoOptionalFields(t *testing.T) {
	p := domain.Person{Name: "Bob", Email: "b@b.com", Document: "11144477735"}
	m := PersonToModel(p)

	if m.Contact != nil {
		t.Error("expected Contact to be nil")
	}
	if m.Address != nil {
		t.Error("expected Address to be nil")
	}
}

func TestPersonToDomain_WithAllFields(t *testing.T) {
	now := time.Now()
	m := personmodel.Model{
		ID:            2,
		Name:          "Bob",
		Email:         "bob@example.com",
		Document:      "11144477735",
		IsActive:      true,
		Contact:       strp("11988887777"),
		Address:       strp("Rua B"),
		AddressNumber: strp("20"),
		Neighborhood:  strp("Vila"),
		City:          strp("RJ"),
		Country:       strp("BR"),
		ZipCode:       strp("20040020"),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	p := PersonToDomain(m)

	if p.ID != 2 || p.Name != "Bob" || p.Email != "bob@example.com" {
		t.Error("basic fields mismatch")
	}
	if p.Contact != "11988887777" {
		t.Error("Contact mismatch")
	}
	if p.Address.Street != "Rua B" {
		t.Error("Address.Street mismatch")
	}
	if p.Address.Number != "20" {
		t.Error("Address.Number mismatch")
	}
	if p.Address.Neighborhood != "Vila" {
		t.Error("Address.Neighborhood mismatch")
	}
	if p.Address.City != *m.City {
		t.Error("Address.City mismatch")
	}
	if p.Address.Country != "BR" {
		t.Error("Address.Country mismatch")
	}
	if p.Address.ZipCode != "20040020" {
		t.Error("Address.ZipCode mismatch")
	}
}

func TestPersonToDomain_WithNilOptionalFields(t *testing.T) {
	m := personmodel.Model{Name: "X", Email: "x@x.com", Document: "52998224725"}
	p := PersonToDomain(m)

	if p.Contact != "" {
		t.Error("expected empty Contact")
	}
	if p.Address.Street != "" {
		t.Error("expected empty Street")
	}
}

// --- User ---

func TestUserToModel(t *testing.T) {
	u := domain.User{ID: 5, Password: "hash", Role: "attendant", PersonID: 3}
	m := UserToModel(u)

	if m.ID != 5 || m.Password != "hash" || m.Role != "attendant" || m.PersonID != 3 {
		t.Error("UserToModel fields mismatch")
	}
}

func TestUserToDomain(t *testing.T) {
	now := time.Now()
	m := usermodel.Model{ID: 5, Password: "hash", Role: "attendant", PersonID: 3, CreatedAt: now, UpdatedAt: now}
	u := UserToDomain(m)

	if u.ID != 5 || u.Password != "hash" || u.Role != "attendant" || u.PersonID != 3 {
		t.Error("UserToDomain fields mismatch")
	}
	if u.CreatedAt != now {
		t.Error("CreatedAt mismatch")
	}
}

// --- Vehicle ---

func TestVehicleToModel(t *testing.T) {
	v := domain.Vehicle{ID: 1, Plate: "ABC1234", Name: "Civic", Year: 2020, Brand: "Honda", Color: "Prata"}
	m := VehicleToModel(v)

	if m.ID != 1 || m.Plate != "ABC1234" || m.Name != "Civic" || m.Year != 2020 || m.Brand != "Honda" || m.Color != "Prata" {
		t.Error("VehicleToModel fields mismatch")
	}
}

func TestVehicleToDomain(t *testing.T) {
	now := time.Now()
	m := vehiclemodel.Model{ID: 1, Plate: "ABC1234", Name: "Civic", Year: 2020, Brand: "Honda", Color: "Prata", CreatedAt: now}
	v := VehicleToDomain(m)

	if v.ID != 1 || v.Plate != "ABC1234" || v.Year != 2020 {
		t.Error("VehicleToDomain fields mismatch")
	}
	if v.CreatedAt != now {
		t.Error("CreatedAt mismatch")
	}
}

// --- Customer ---

func TestCustomerToModel(t *testing.T) {
	c := domain.Customer{ID: 2, Type: "pf", PersonID: 7}
	m := CustomerToModel(c)

	if m.ID != 2 || m.Type != "pf" || m.PersonID != 7 {
		t.Error("CustomerToModel fields mismatch")
	}
}

func TestCustomerToDomain(t *testing.T) {
	now := time.Now()
	m := customermodel.Model{ID: 2, Type: "pf", PersonID: 7, CreatedAt: now}
	c := CustomerToDomain(m)

	if c.ID != 2 || c.Type != "pf" || c.PersonID != 7 {
		t.Error("CustomerToDomain fields mismatch")
	}
}

// --- Employee ---

func TestEmployeeToModel(t *testing.T) {
	e := domain.Employee{ID: 3, Position: "Mechanic", PersonID: 5}
	m := EmployeeToModel(e)

	if m.ID != 3 || m.Position != "Mechanic" || m.PersonID != 5 {
		t.Error("EmployeeToModel fields mismatch")
	}
}

func TestEmployeeToDomain(t *testing.T) {
	now := time.Now()
	m := employeemodel.Model{ID: 3, Position: "Mechanic", PersonID: 5, CreatedAt: now}
	e := EmployeeToDomain(m)

	if e.ID != 3 || e.Position != "Mechanic" || e.PersonID != 5 {
		t.Error("EmployeeToDomain fields mismatch")
	}
}

// --- Company ---

func TestCompanyToModel_WithAddress(t *testing.T) {
	c := domain.Company{
		ID: 1, Name: "ACME", Email: "acme@acme.com",
		Document: "11222333000181", Contact: "11999999999",
		Address: domain.Address{Street: "Av. Paulista", Number: "1000", City: "SP", Country: "BR", ZipCode: "01310100"},
	}
	m := CompanyToModel(c)

	if m.ID != 1 || m.Name != "ACME" {
		t.Error("CompanyToModel basic fields mismatch")
	}
	if m.Address == nil || *m.Address != "Av. Paulista" {
		t.Error("expected Address to be set")
	}
	if m.City == nil || *m.City != "SP" {
		t.Error("expected City to be set")
	}
}

func TestCompanyToModel_WithoutAddress(t *testing.T) {
	c := domain.Company{Name: "X", Email: "x@x.com", Document: "11222333000181", Contact: "11"}
	m := CompanyToModel(c)

	if m.Address != nil {
		t.Error("expected Address to be nil")
	}
}

func TestCompanyToDomain_WithAddress(t *testing.T) {
	now := time.Now()
	m := companymodel.Model{
		ID: 1, Name: "ACME", Email: "acme@acme.com", Document: "11222333000181", Contact: "11999999999",
		Address: strp("Av. Paulista"), AddressNumber: strp("1000"),
		Neighborhood: strp("Bela Vista"), City: strp("SP"), Country: strp("BR"), ZipCode: strp("01310100"),
		CreatedAt: now,
	}
	c := CompanyToDomain(m)

	if c.Address.Street != "Av. Paulista" || c.Address.Number != "1000" {
		t.Error("CompanyToDomain address mismatch")
	}
	if c.Address.Neighborhood != "Bela Vista" || c.Address.City != "SP" {
		t.Error("CompanyToDomain neighborhood/city mismatch")
	}
}

func TestCompanyToDomain_WithoutAddress(t *testing.T) {
	m := companymodel.Model{Name: "X", Email: "x@x.com", Document: "11222333000181", Contact: "11"}
	c := CompanyToDomain(m)

	if c.Address.Street != "" {
		t.Error("expected empty Street")
	}
}

// --- CustomerVehicle ---

func TestCustomerVehicleToModel(t *testing.T) {
	cv := domain.CustomerVehicle{ID: 1, CustomerID: 2, VehicleID: 3}
	m := CustomerVehicleToModel(cv)

	if m.ID != 1 || m.CustomerID != 2 || m.VehicleID != 3 {
		t.Error("CustomerVehicleToModel fields mismatch")
	}
}

func TestCustomerVehicleToDomain(t *testing.T) {
	now := time.Now()
	m := cvmodel.Model{
		ID: 1, CustomerID: 2, VehicleID: 3,
		Vehicle:   vehiclemodel.Model{ID: 3, Plate: "ABC1234", Name: "Civic", Year: 2020, Brand: "Honda", Color: "Prata"},
		CreatedAt: now,
	}
	cv := CustomerVehicleToDomain(m)

	if cv.ID != 1 || cv.CustomerID != 2 || cv.VehicleID != 3 {
		t.Error("CustomerVehicleToDomain fields mismatch")
	}
	if cv.Vehicle.Plate != "ABC1234" {
		t.Error("expected Vehicle.Plate to be set")
	}
}
