package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OwnerID   primitive.ObjectID `json:"ownerId" bson:"ownerId"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}
