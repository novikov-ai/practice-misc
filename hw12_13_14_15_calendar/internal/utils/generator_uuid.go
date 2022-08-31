package utils

import "github.com/google/uuid"

func GenerateUUID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
