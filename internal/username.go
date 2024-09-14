package internal

import (
	"errors"
	"regexp"
)

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	re := regexp.MustCompile("^[a-zA-Z0-9_]+$")
	if !re.MatchString(username) {
		return errors.New("username can only contain alphanumeric characters and underscores")
	}

	return nil
}
