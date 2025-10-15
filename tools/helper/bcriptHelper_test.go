package helper_test

import (
	"github.com/nurefendi/auth-service-using-golang/tools/helper"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashBcript(t *testing.T) {
	password := "mypassword"
	hash, err := helper.HashBcript(password)
	assert.NoError(t, err, "Hashing should not return an error")
	assert.NotEmpty(t, hash, "Generated hash should not be empty")
}

func TestCompareHashBcript(t *testing.T) {
	password := "mypassword"
	hash, err := helper.HashBcript(password)
	assert.NoError(t, err, "Hashing should not return an error")

	err = helper.CompareHashBcript(password, hash)
	assert.NoError(t, err, "Password should match the hash")
}

func TestCompareHashBcriptInvalidPassword(t *testing.T) {
	password := "mypassword"
	wrongPassword := "wrongpassword"
	hash, err := helper.HashBcript(password)
	assert.NoError(t, err, "Hashing should not return an error")

	err = helper.CompareHashBcript(wrongPassword, hash)
	assert.Error(t, err, "Wrong password should not match the hash")
	assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err, "Error should be ErrMismatchedHashAndPassword")
}
