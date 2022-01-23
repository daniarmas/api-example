package repository

import "golang.org/x/crypto/bcrypt"

type HashPasswordQuery interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type hashPasswordQuery struct{}

func (i *hashPasswordQuery) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func (i *hashPasswordQuery) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
