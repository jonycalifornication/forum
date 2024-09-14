package internal

import "github.com/gofrs/uuid"

func GenerateSessionID() (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}
