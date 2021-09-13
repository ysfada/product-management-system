package hasher

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/ysfada/product-management-system/util"
	"golang.org/x/crypto/argon2"
)

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func (p *Argon2Params) Defaults() {
	p.Memory = 64 * 1024
	p.Iterations = 1
	p.Parallelism = 2
	p.SaltLength = 16
	p.KeyLength = 32
}

type Argon2 struct {
	Params *Argon2Params
}

var _ Hasher = (*Argon2)(nil)

func NewArgon2(argon2Params ...*Argon2Params) Hasher {
	hasher := &Argon2{}

	if argon2Params == nil {
		hasher.Params = &Argon2Params{}
		hasher.Params.Defaults()
	} else if len(argon2Params) == 1 {
		hasher.Params = argon2Params[0]
	} else {
		panic("len(params) must not be longer than 1")
	}

	return hasher
}

func (h *Argon2) Hash(password string) ([]byte, error) {
	salt, err := util.GenerateRandomBytes(h.Params.SaltLength)
	if err != nil {
		return nil, err
	}

	key := argon2.IDKey([]byte(password), salt, h.Params.Iterations, h.Params.Memory, h.Params.Parallelism, h.Params.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	hash := []byte(fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, h.Params.Memory, h.Params.Iterations, h.Params.Parallelism, b64Salt, b64Key))
	return hash, nil
}

func (h *Argon2) Compare(hashedPassword, password string) error {
	params, salt, key, err := h.decodeHash(hashedPassword)
	if err != nil {
		return err
	}

	otherKey := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	keyLen := int32(len(key))
	otherKeyLen := int32(len(otherKey))

	if subtle.ConstantTimeEq(keyLen, otherKeyLen) == 0 {
		return errors.New("argon2id: hashedPassword is not the hash of the given password")
	}
	if subtle.ConstantTimeCompare(key, otherKey) == 1 {
		return nil
	}
	return errors.New("argon2id: hashedPassword is not the hash of the given password")
}

func (h *Argon2) decodeHash(hash string) (params *Argon2Params, salt, key []byte, err error) {
	vals := strings.Split(hash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("argon2id: hash is not in the correct format")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("argon2id: incompatible version of argon2")
	}

	params = &Argon2Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = uint32(len(salt))

	key, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = uint32(len(key))

	return params, salt, key, nil
}
