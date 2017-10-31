package dependencies_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t-ksn/user-service/src/dependencies"
)

func TestPasswordHasher_HashThenVerify(t *testing.T) {
	pwdHasher := dependencies.MakePasswordHasher()
	pwd := "my secret password"

	pwdHash, err := pwdHasher.Hash(pwd)
	assert.NoError(t, err)

	verified := pwdHasher.Verify(pwd, pwdHash)
	assert.True(t, verified)
}
