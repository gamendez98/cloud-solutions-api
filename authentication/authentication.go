package authentication

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JwtCustomClaims represents the custom claims structure for JWT including a username and standard registered claims.
type JwtCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// HashPassword hashes the given password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a hashed password with its plaintext equivalent.
func CheckPasswordHash(password, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}
