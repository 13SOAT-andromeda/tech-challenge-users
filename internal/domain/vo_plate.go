package domain

import (
	"errors"
	"regexp"
	"strings"
)

type Plate string

var (
	rePlateOld      = regexp.MustCompile(`^[A-Z]{3}\d{4}$`)
	rePlateMercosul = regexp.MustCompile(`^[A-Z]{3}\d[A-Z]\d{2}$`)
)

func NewPlate(raw string) (Plate, error) {
	normalized := strings.ToUpper(strings.ReplaceAll(raw, "-", ""))
	if !rePlateOld.MatchString(normalized) && !rePlateMercosul.MatchString(normalized) {
		return "", errors.New("invalid vehicle plate format")
	}
	return Plate(normalized), nil
}

func (p Plate) String() string { return string(p) }
