package domain

import "testing"

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "Admin@123", false},
		{"too short", "Ab@1", true},
		{"no uppercase", "admin@123", true},
		{"no lowercase", "ADMIN@123", true},
		{"no digit", "Admin@abc", true},
		{"no special", "Admin1234", true},
		{"exactly 8 chars valid", "Abcde@1!", false},
		{"only spaces", "        ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPassword(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPassword(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
