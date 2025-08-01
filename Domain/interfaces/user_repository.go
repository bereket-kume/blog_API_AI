package interfaces 



import "blog-api/Domain/models"

type UserRepository interface {
	Insert(user models.User) error
	FindByEmail(email string) (models.User, error)
	UpdatePass(user models.User) error
	UpdateRole(user models.User) error
	Delete(email string) error
	Verify(email string) error
	CountUsers() (int64, error) 
}
