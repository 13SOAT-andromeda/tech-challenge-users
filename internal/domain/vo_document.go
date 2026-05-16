package domain

import (
	"errors"
	"regexp"
)

type Document string

var nonDigit = regexp.MustCompile(`\D`)

func NewDocument(raw string) (Document, error) {
	digits := nonDigit.ReplaceAllString(raw, "")
	switch len(digits) {
	case 11:
		if !validCPF(digits) {
			return "", errors.New("invalid CPF")
		}
	case 14:
		if !validCNPJ(digits) {
			return "", errors.New("invalid CNPJ")
		}
	default:
		return "", errors.New("document must be a valid CPF (11 digits) or CNPJ (14 digits)")
	}
	return Document(digits), nil
}

func (d Document) String() string { return string(d) }

func validCPF(s string) bool {
	if allSameDigits(s) {
		return false
	}
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(s[i]-'0') * (10 - i)
	}
	r := (sum * 10) % 11
	if r >= 10 {
		r = 0
	}
	if r != int(s[9]-'0') {
		return false
	}
	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(s[i]-'0') * (11 - i)
	}
	r = (sum * 10) % 11
	if r >= 10 {
		r = 0
	}
	return r == int(s[10]-'0')
}

func validCNPJ(s string) bool {
	if allSameDigits(s) {
		return false
	}
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 12; i++ {
		sum += int(s[i]-'0') * weights1[i]
	}
	r := sum % 11
	if r < 2 {
		r = 0
	} else {
		r = 11 - r
	}
	if r != int(s[12]-'0') {
		return false
	}
	sum = 0
	for i := 0; i < 13; i++ {
		sum += int(s[i]-'0') * weights2[i]
	}
	r = sum % 11
	if r < 2 {
		r = 0
	} else {
		r = 11 - r
	}
	return r == int(s[13]-'0')
}

func allSameDigits(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}
