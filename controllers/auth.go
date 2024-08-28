package controllers

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"go_react_app/models"
	"go_react_app/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Email  string             `json:"email"`
	UserID primitive.ObjectID `json:"userId"`
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

func Register(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid input"})
	}

	if user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Email and password cannot be empty",
		})
	}

	var existingUser models.User
	err := utils.GetUserCollection().FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already in use"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not hash password"})
	}
	user.Password = string(hashedPassword)

	_, err = utils.GetUserCollection().InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User registered successfully"})
}

func Login(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid input"})
	}

	if user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Email and password cannot be empty",
		})
	}

	var existingUser models.User
	err := utils.GetUserCollection().FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	expirationTime := time.Now().Add(24 * time.Hour) // O token expira em 24 horas
	claims := &Claims{
		Email:  user.Email,
		UserID: existingUser.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create token"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Login successful", "token": tokenString, "user": existingUser.Name})
}

func getAuthenticatedUserID(c *fiber.Ctx) (primitive.ObjectID, error) {
	tokenString := c.Get("Authorization")

	if tokenString == "" {
		return primitive.NilObjectID, errors.New("missing token")
	}

	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !tkn.Valid {
		return primitive.NilObjectID, errors.New("invalid token")
	}

	return claims.UserID, nil
}
