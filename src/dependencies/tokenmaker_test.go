package dependencies_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t-ksn/core-kit/apierror"
	"github.com/t-ksn/user-service/src/dependencies"
	"github.com/t-ksn/user-service/src/service"
)

func TestTokenMaker_MakeThenVerify(t *testing.T) {
	expected := service.Token{UserID: "123"}
	tm := dependencies.MakeTokenGenerator("my secret")
	tokenStr, err := tm.Make(expected)
	assert.NoError(t, err)

	token, err := tm.Verify(tokenStr)
	assert.NoError(t, err)
	assert.Equal(t, expected, token)
}

func TestTokenMaker_Verify_SendIncorrectString_ReturnErr(t *testing.T) {
	tm := dependencies.MakeTokenGenerator("my secret")
	_, err := tm.Verify("incorrect string")

	assert.Equal(t, apierror.UnauthorizedRequestErr, err)
}

func TestTokenMaker_Verify_SendTokenWithOtherSecretSign_ReturnErr(t *testing.T) {
	tm := dependencies.MakeTokenGenerator("my secret")
	token, err := tm.Make(service.Token{UserID: "123"})

	newTm := dependencies.MakeTokenGenerator("my new secret")
	_, err = newTm.Verify(token)

	assert.Equal(t, apierror.UnauthorizedRequestErr, err)
}

func TestTokenMaker_Verify_a_ReturnErr(t *testing.T) {
	tm := dependencies.MakeTokenGenerator("my secret")
	_, err := tm.Verify("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")

	assert.Equal(t, apierror.UnauthorizedRequestErr, err)
}
