package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email   string             `bson:"email,omitempty" json:"email,omitempty"`
	Bio     string             `bson:"bio,omitempty" json:"bio,omitempty"`
	Picture string             `bson:"picture,omitempty" json:"picture,omitempty"`
	Contact string             `bson:"contact,omitempty" json:"contact,omitempty"`
}
