package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(input string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(input), 14)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
