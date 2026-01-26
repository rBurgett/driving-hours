package utils

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail checks if an email address is valid
func ValidateEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}

// ValidatePassword checks if a password meets minimum requirements
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}
	return true, ""
}

// SanitizeString trims whitespace and removes control characters
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// ValidateName checks if a name is valid
func ValidateName(name string) (bool, string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return false, "Name is required"
	}
	if len(name) > 100 {
		return false, "Name must be less than 100 characters"
	}
	return true, ""
}

// ValidateDate checks if a date string is in YYYY-MM-DD format
func ValidateDate(date string) bool {
	if len(date) != 10 {
		return false
	}
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return dateRegex.MatchString(date)
}

// ValidateHours checks if hours value is reasonable
func ValidateHours(hours float64) (bool, string) {
	if hours < 0 {
		return false, "Hours cannot be negative"
	}
	if hours > 24 {
		return false, "Hours cannot exceed 24"
	}
	return true, ""
}
