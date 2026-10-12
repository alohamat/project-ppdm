package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID    bson.ObjectID `bson:"_id,omitempty"`
	Email string `json:"email"`
	Password string `json:"password"`
}