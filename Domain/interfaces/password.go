package interfaces

type Hasher interface {
	HashPassword(password string) string
	VerifyPassword(hashed, password string) bool
}
