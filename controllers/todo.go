package controllers

import (
	"context"
	"go_react_app/models"
	"go_react_app/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTodos(c *fiber.Ctx) error {
	ownerID, err := getAuthenticatedUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var todos []models.Todo
	filter := bson.M{"ownerId": ownerID}

	cursor, err := utils.GetTodoCollection().Find(context.Background(), filter)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo models.Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func CreateTodos(c *fiber.Ctx) error {
	todo := new(models.Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Todo body cannot be empty",
		})
	}

	ownerID, err := getAuthenticatedUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	todo.OwnerID = ownerID

	insertResult, err := utils.GetTodoCollection().InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(fiber.Map{"message": "Todo created successfully", "todo": todo})
}

func UpdateTodos(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid todo ID",
		})
	}

	ownerID, err := getAuthenticatedUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	filter := bson.M{"_id": objectID, "ownerId": ownerID}
	update := bson.M{"$set": bson.M{"completed": true}}

	result, err := utils.GetTodoCollection().UpdateOne(context.Background(), filter, update)
	if err != nil || result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found or you do not have permission to update"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Todo updated successfully",
	})
}

func DeleteTodos(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid todo ID",
		})
	}

	ownerID, err := getAuthenticatedUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	filter := bson.M{"_id": objectID, "ownerId": ownerID}

	result, err := utils.GetTodoCollection().DeleteOne(context.Background(), filter)
	if err != nil || result.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found or you do not have permission to delete"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Todo deleted successfully",
	})
}
