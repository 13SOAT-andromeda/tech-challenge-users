package domain

import (
	"errors"
	"unicode"
)

type Password string

func NewPassword(plain string) (Password, error) {
	if len(plain) < 8 {
		return "", errors.New("password must be at least 8 characters")
	}
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range plain {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}
	if !hasUpper {
		return "", errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return "", errors.New("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return "", errors.New("password must contain at least one digit")
	}
	if !hasSpecial {
		return "", errors.New("password must contain at least one special character")
	}
	return Password(plain), nil
}

func (p Password) String() string { return string(p) }
