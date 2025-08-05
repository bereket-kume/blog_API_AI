package interfaces

import "blog-api/Domain/models"

type UserRepository interface {
	Insert(user models.User) error
	FindByEmail(email string) (models.User, error)
	UpdatePass(email string, passowrdHash string) error
	UpdateRole(email, role string) error
	Delete(email string) error
	Verify(email string) error
	CountUsers() (int64, error)
}
