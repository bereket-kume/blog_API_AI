package interfaces

type Hasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashed, password string) bool
}
