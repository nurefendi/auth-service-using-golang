package helper

import "golang.org/x/crypto/bcrypt"

func HashBcript(text string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(text), 14)
	return string(bytes), err
}

func CompareHashBcript(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}