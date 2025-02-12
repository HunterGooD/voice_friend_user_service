package utils

import (
	"regexp"
	"strings"
)

var (
	reEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	rePhone = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

// ValidateEmail regexp validation
// if match eq false it's not correct
func ValidateEmail(email string) bool {
	return reEmail.MatchString(email)
}

func ValidatePhone(phone string) bool {
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	return rePhone.MatchString(cleaned)
}
