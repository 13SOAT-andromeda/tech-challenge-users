package domain

import "testing"

func TestNewPlate_OldFormat(t *testing.T) {
	p, err := NewPlate("ABC1234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.String() != "ABC1234" {
		t.Errorf("expected 'ABC1234', got %q", p.String())
	}
}

func TestNewPlate_OldFormatWithDash(t *testing.T) {
	p, err := NewPlate("ABC-1234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.String() != "ABC1234" {
		t.Errorf("expected 'ABC1234' after normalization, got %q", p.String())
	}
}

func TestNewPlate_MercosulFormat(t *testing.T) {
	p, err := NewPlate("ABC1D23")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.String() != "ABC1D23" {
		t.Errorf("expected 'ABC1D23', got %q", p.String())
	}
}

func TestNewPlate_LowercaseNormalized(t *testing.T) {
	p, err := NewPlate("abc1234")
	if err != nil {
		t.Fatalf("expected lowercase plate to be normalized: %v", err)
	}
	if p.String() != "ABC1234" {
		t.Errorf("expected 'ABC1234', got %q", p.String())
	}
}

func TestNewPlate_Invalid(t *testing.T) {
	cases := []string{"INVALID", "1234567", "AB12345", "ABC123", "ABCD1234", ""}
	for _, raw := range cases {
		_, err := NewPlate(raw)
		if err == nil {
			t.Errorf("expected error for invalid plate %q", raw)
		}
	}
}

func TestPlate_String(t *testing.T) {
	p := Plate("ABC1234")
	if p.String() != "ABC1234" {
		t.Errorf("expected 'ABC1234', got %q", p.String())
	}
}
