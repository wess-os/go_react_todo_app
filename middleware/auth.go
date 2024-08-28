package middleware

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

func IsAuthenticated(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	if tokenString == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing token"})
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	c.Locals("email", claims.Email)

	return c.Next()
}
