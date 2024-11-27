package hashing_test

import (
	"testing"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/hashing"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Run("Successful Hashing", func(t *testing.T) {
		password := "securePassword123!"
		hash, err := hashing.HashPassword(password)
		require.NoError(t, err, "Hashing password should not return an error")
		require.NotEmpty(t, hash, "Hash should not be empty")
		require.NotEqual(t, password, hash, "Hash should not be equal to the original password")
	})

	t.Run("Empty Password", func(t *testing.T) {
		password := ""
		hash, err := hashing.HashPassword(password)
		require.NoError(t, err, "Hashing an empty password should not return an error")
		require.NotEmpty(t, hash, "Hash should not be empty even for an empty password")
	})
}

func TestCheckPasswordHash(t *testing.T) {
	t.Run("Valid Password and Hash", func(t *testing.T) {
		password := "securePassword123!"
		hash, err := hashing.HashPassword(password)
		require.NoError(t, err, "Hashing password should not return an error")

		valid := hashing.CheckPasswordHash(password, hash)
		require.True(t, valid, "Password should match the hash")
	})

	t.Run("Invalid Password", func(t *testing.T) {
		password := "securePassword123!"
		hash, err := hashing.HashPassword(password)
		require.NoError(t, err, "Hashing password should not return an error")

		invalidPassword := "wrongPassword123!"
		valid := hashing.CheckPasswordHash(invalidPassword, hash)
		require.False(t, valid, "Wrong password should not match the hash")
	})

	t.Run("Empty Password and Valid Hash", func(t *testing.T) {
		password := ""
		hash, err := hashing.HashPassword(password)
		require.NoError(t, err, "Hashing empty password should not return an error")

		valid := hashing.CheckPasswordHash(password, hash)
		require.True(t, valid, "Empty password should match the hash of an empty password")
	})

	t.Run("Empty Hash", func(t *testing.T) {
		password := "securePassword123!"
		valid := hashing.CheckPasswordHash(password, "")
		require.False(t, valid, "Password should not match an empty hash")
	})
}
