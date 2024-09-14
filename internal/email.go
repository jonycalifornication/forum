package internal

import (
	"errors"
	"net"
	"regexp"
	"strings"
)

// ValidateEmail checks if the email address is in a valid format and if the domain exists.
func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}

	// Extract domain from email
	domain := email[strings.LastIndex(email, "@")+1:]

	// Check if domain has MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return errors.New("email domain does not exist")
	}

	return nil
}
