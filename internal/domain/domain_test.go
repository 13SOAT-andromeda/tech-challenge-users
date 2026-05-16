package domain

import "testing"

func TestValidateCustomerType_Valid(t *testing.T) {
	if err := ValidateCustomerType(CustomerTypePF); err != nil {
		t.Errorf("expected 'pf' to be valid: %v", err)
	}
	if err := ValidateCustomerType(CustomerTypePJ); err != nil {
		t.Errorf("expected 'pj' to be valid: %v", err)
	}
}

func TestValidateCustomerType_Invalid(t *testing.T) {
	if err := ValidateCustomerType("corporate"); err == nil {
		t.Error("expected error for invalid customer type")
	}
	if err := ValidateCustomerType(""); err == nil {
		t.Error("expected error for empty customer type")
	}
}

func TestValidateRole_Valid(t *testing.T) {
	roles := []string{RoleCustomer, RoleAttendant, RoleMechanic, RoleAdministrator}
	for _, r := range roles {
		if err := ValidateRole(r); err != nil {
			t.Errorf("expected role %q to be valid: %v", r, err)
		}
	}
}

func TestValidateRole_Invalid(t *testing.T) {
	if err := ValidateRole("superadmin"); err == nil {
		t.Error("expected error for invalid role")
	}
	if err := ValidateRole(""); err == nil {
		t.Error("expected error for empty role")
	}
}

func TestPassword_String(t *testing.T) {
	p, err := NewPassword("Secure@123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.String() != "Secure@123" {
		t.Errorf("expected 'Secure@123', got %q", p.String())
	}
}
