package hashing

import (
	"fmt"

	bcr "golang.org/x/crypto/bcrypt"
)

type bcrypt struct{}

var _ Hasher = (*bcrypt)(nil)

// NewBcrypt is used to create an instance of bcrypt hasher.
func NewBcrypt() *bcrypt {
	return &bcrypt{}
}

func (b *bcrypt) GenerateHashFromString(value string) (string, error) {
	hashedValue, err := bcr.GenerateFromPassword([]byte(value), 14)
	if err != nil {
		return "", fmt.Errorf("hash value: %w", err)
	}

	return string(hashedValue), nil
}
