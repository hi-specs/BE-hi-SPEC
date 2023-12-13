package enkrip

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type HashInterface interface {
	Compare(hashed string, input string) error
	HashPassword(input string) (string, error)
}

type hash struct{}

func New() HashInterface {
	return &hash{}
}

func (h *hash) Compare(hashed string, input string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input))
}

func (h *hash) HashPassword(input string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		log.Println("HASH - terjadi kesalahan saat hash password, error", err.Error())
		return "", err
	}

	return string(result), nil
}
