package authentication

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
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

func GetCurrentUsername(c echo.Context) (string, error) {
	user := c.Get("user")
	token, ok := user.(*jwt.Token)
	if !ok {
		return "", fmt.Errorf("failed to get JWT token")
	}

	claims := token.Claims
	fmt.Println(claims)
	parsedClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to parse claims")
	}
	username, ok := parsedClaims["username"].(string)
	if !ok {
		return "", fmt.Errorf("username not found in token claims")
	}

	return username, nil
}
