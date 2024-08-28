package routes

import (
	"go_react_app/controllers"
	"go_react_app/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/api/todos", middleware.IsAuthenticated, controllers.GetTodos)
	app.Post("/api/todos", middleware.IsAuthenticated, controllers.CreateTodos)
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Patch("/api/todos/:id", middleware.IsAuthenticated, controllers.UpdateTodos)
	app.Delete("/api/todos/:id", middleware.IsAuthenticated, controllers.DeleteTodos)
}
