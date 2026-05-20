package domain

import "testing"

func TestNewDocument_ValidCPFWithMask(t *testing.T) {
	doc, err := NewDocument("529.982.247-25")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc.String() != "52998224725" {
		t.Errorf("expected '52998224725', got %q", doc.String())
	}
}

func TestNewDocument_ValidCPFNoMask(t *testing.T) {
	_, err := NewDocument("11144477735")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewDocument_ValidCNPJWithMask(t *testing.T) {
	doc, err := NewDocument("11.222.333/0001-81")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc.String() != "11222333000181" {
		t.Errorf("expected '11222333000181', got %q", doc.String())
	}
}

func TestNewDocument_ValidCNPJNoMask(t *testing.T) {
	_, err := NewDocument("11222333000181")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewDocument_AllSameDigitsCPF(t *testing.T) {
	_, err := NewDocument("111.111.111-11")
	if err == nil {
		t.Error("expected error for CPF with all same digits")
	}
}

func TestNewDocument_InvalidCPFFirstDigit(t *testing.T) {
	_, err := NewDocument("529.982.247-00")
	if err == nil {
		t.Error("expected error for CPF with wrong first check digit")
	}
}

func TestNewDocument_InvalidCPFSecondDigit(t *testing.T) {
	_, err := NewDocument("529.982.247-26")
	if err == nil {
		t.Error("expected error for CPF with wrong second check digit")
	}
}

func TestNewDocument_AllSameDigitsCNPJ(t *testing.T) {
	_, err := NewDocument("11.111.111/1111-11")
	if err == nil {
		t.Error("expected error for CNPJ with all same digits")
	}
}

func TestNewDocument_InvalidCNPJFirstCheckDigit(t *testing.T) {
	_, err := NewDocument("11222333000191")
	if err == nil {
		t.Error("expected error for CNPJ with wrong first check digit")
	}
}

func TestNewDocument_InvalidCNPJSecondCheckDigit(t *testing.T) {
	_, err := NewDocument("11222333000182")
	if err == nil {
		t.Error("expected error for CNPJ with wrong second check digit")
	}
}

func TestNewDocument_TooShort(t *testing.T) {
	_, err := NewDocument("12345")
	if err == nil {
		t.Error("expected error for document with wrong digit count")
	}
}

func TestNewDocument_TooLong(t *testing.T) {
	_, err := NewDocument("123456789012345")
	if err == nil {
		t.Error("expected error for document with too many digits")
	}
}

func TestDocument_String(t *testing.T) {
	doc := Document("52998224725")
	if doc.String() != "52998224725" {
		t.Errorf("expected '52998224725', got %q", doc.String())
	}
}
