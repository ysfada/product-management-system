package hasher

type Hasher interface {
	Hash(password string) ([]byte, error)
	Compare(hashedPassword, password string) error
}
