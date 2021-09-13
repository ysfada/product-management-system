package hasher

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct {
	Cost int
}

var _ Hasher = (*Bcrypt)(nil)

func NewBcrypt(cost ...int) Hasher {
	hasher := &Bcrypt{}
	if cost == nil {
		hasher.Cost = bcrypt.DefaultCost
	} else if len(cost) == 1 {
		hasher.Cost = cost[0]
	} else {
		panic("len(cost) must not be longer than 1")
	}

	return hasher
}

func (h *Bcrypt) Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), h.Cost)
}

func (h *Bcrypt) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
