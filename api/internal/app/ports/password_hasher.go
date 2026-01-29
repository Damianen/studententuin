package ports

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(passwrod string, hash string) (bool, error)
}
