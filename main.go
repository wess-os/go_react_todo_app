package main

import (
	"log"
	"os"

	"go_react_app/routes"
	"go_react_app/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file: ", err)
		}
	}

	app := fiber.New()

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: os.Getenv("ROTA_FRONTEND"),
	// 	AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	// }))

	if os.Getenv("ENV") == "production" {
		app.Static("/", "./client/dist")
	}

	utils.ConnectDB() // Conectar ao banco de dados

	defer utils.DisconnectDB() // Desconectar ao final

	routes.SetupRoutes(app) // Configurar as rotas

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}
