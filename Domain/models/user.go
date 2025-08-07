package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// in Domain/models/user.go
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password,omitempty" json:"-"`
	Role     string             `bson:"role" json:"role"`
	Verified bool               `bson:"verified" json:"verified"`

	Bio     string `bson:"bio,omitempty" json:"bio,omitempty"`
	Picture string `bson:"picture,omitempty" json:"picture,omitempty"`
	Contact string `bson:"contact,omitempty" json:"contact,omitempty"`
}
