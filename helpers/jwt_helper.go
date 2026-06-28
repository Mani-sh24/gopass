package helpers

import (
	"log"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func InitJWT(secret string) error {
	if secret == "" {
		log.Fatal("NO JWT SECRET FOUND")
	}

	jwtSecret = []byte(secret)
	return nil
}

func GenerateJWT(id, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"id":    id,
	})

	// secret := os.Getenv("JWT_TOK")

	return token.SignedString([]byte(jwtSecret))
}
