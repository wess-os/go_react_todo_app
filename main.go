package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var collection *mongo.Collection
var userCollection *mongo.Collection
var jwtKey = []byte(os.Getenv("SECRET_KEY"))

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	collection = client.Database("golang_db").Collection("todos")
	userCollection = client.Database("golang_db").Collection("users")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5174",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/api/todos", isAuthenticated, getTodos)
	app.Post("/api/todos", isAuthenticated, createTodos)
	app.Post("/api/register", register)
	app.Post("/api/login", login)
	app.Patch("/api/todos/:id", isAuthenticated, updateTodos)
	app.Delete("/api/todos/:id", isAuthenticated, deleteTodos)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	fmt.Println("Listening on port " + port)

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func createTodos(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Todo body cannot be empty",
		})
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}

func updateTodos(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid todo ID",
		})
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Todo updated successfully",
	})
}

func deleteTodos(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid todo ID",
		})
	}

	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Todo deleted successfully",
	})
}

func register(c *fiber.Ctx) error {
	var user User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid input"})
	}

	if user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Email and password cannot be empty",
		})
	}

	var existingUser User
	err := userCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already in use"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not hash password"})
	}
	user.Password = string(hashedPassword)

	_, err = userCollection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User registered successfully"})
}

func login(c *fiber.Ctx) error {
	var user User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid input"})
	}

	var existingUser User
	err := userCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	expirationTime := time.Now().Add(24 * time.Hour) // O token expira em 24 horas
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create token"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Login successful", "token": tokenString})
}

func isAuthenticated(c *fiber.Ctx) error {
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
