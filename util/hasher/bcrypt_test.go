package hasher

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptHash(t *testing.T) {
	bcrypt := NewBcrypt()

	hash1, err := bcrypt.Hash("pa$$word")
	assert.Nil(t, err)

	hash2, err := bcrypt.Hash("pa$$word")
	assert.Nil(t, err)

	assert.NotEqual(t, 0, strings.Compare(string(hash1), string(hash2)), "hashes must be unique")
}

func TestBcryptCompare(t *testing.T) {
	bcrypt := NewBcrypt()

	hash, err := bcrypt.Hash("pa$$word")
	assert.Nil(t, err)

	err = bcrypt.Compare("pa$$word", string(hash))
	assert.NotNil(t, err)

	err = bcrypt.Compare("otherPa$$word", string(hash))
	assert.NotNil(t, err)
}
