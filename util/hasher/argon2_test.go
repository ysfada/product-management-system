package hasher

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgon2Hash(t *testing.T) {
	hashRX, err := regexp.Compile(`^\$argon2id\$v=19\$m=65536,t=1,p=2\$[A-Za-z0-9+/]{22}\$[A-Za-z0-9+/]{43}$`)
	assert.Nil(t, err)

	argon2 := NewArgon2()

	hash1, err := argon2.Hash("pa$$word")

	assert.Nil(t, err)

	assert.Truef(t, hashRX.MatchString(string(hash1)), "hash %q not in correct format", hash1)

	hash2, err := argon2.Hash("pa$$word")
	assert.Nil(t, err)

	assert.NotEqual(t, 0, strings.Compare(string(hash1), string(hash2)), "hashes must be unique")
}

func TestArgon2Compare(t *testing.T) {
	argon2 := NewArgon2()

	hash, err := argon2.Hash("pa$$word")
	assert.Nil(t, err)

	err = argon2.Compare("pa$$word", string(hash))
	assert.NotNil(t, err)

	err = argon2.Compare("otherPa$$word", string(hash))
	assert.NotNil(t, err)
}
