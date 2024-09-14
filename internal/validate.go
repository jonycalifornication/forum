package internal

import (
	"errors"
	"html"
	"unicode"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var hasLetter, hasNumber, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsLetter(ch):
			hasLetter = true
		case unicode.IsNumber(ch):
			hasNumber = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if !hasLetter || !hasNumber || !hasSpecial {
		return errors.New("password must contain at least one letter, one number, and one special character")
	}

	return nil
}

func SanitizeComment(commentText string) string {
	return html.EscapeString(commentText)
}
