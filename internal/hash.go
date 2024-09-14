package internal

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	return hash, nil
}

func CheckPasswordHash(password string, hash []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	return nil
}
