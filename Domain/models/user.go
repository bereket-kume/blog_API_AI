package models

// User represents the domain model for a user
// Following clean architecture: domain layer should be independent of infrastructure
type User struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password,omitempty" json:"-"`
	Role     string `bson:"role" json:"role"`
	Verified bool   `bson:"verified" json:"verified"`

	Bio     string `bson:"bio,omitempty" json:"bio,omitempty"`
	Picture string `bson:"picture,omitempty" json:"picture,omitempty"`
	Contact string `bson:"contact,omitempty" json:"contact,omitempty"`
}
