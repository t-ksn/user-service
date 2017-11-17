package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/t-ksn/core-kit/apierror"
	"github.com/t-ksn/user-service/src/service"
)

func TestService_Register_Seccess_ReturnNil(t *testing.T) {
	serviceEnv := makeServiceEnv()
	callEnv := setupServcieRegisterEnv(serviceEnv)

	err := serviceEnv.Service.Register(context.Background(), callEnv.Request)

	assert.NoError(t, err)
}

func TestService_Register_PasswordLessThen4Chars_ReturnErrPasswordLessThen4Chars(t *testing.T) {
	serviceEnv := makeServiceEnv()
	callEnv := setupServcieRegisterEnv(serviceEnv)
	callEnv.Request.Password = "123"

	err := serviceEnv.Service.Register(context.Background(), callEnv.Request)

	assert.Equal(t, service.ErrPasswordLessThen4Chars, err)
}

func TestService_Register_UserNameIsEmpty_ReturnErrUserNameIsEmpty(t *testing.T) {
	serviceEnv := makeServiceEnv()
	callEnv := setupServcieRegisterEnv(serviceEnv)
	callEnv.Request.Name = ""

	err := serviceEnv.Service.Register(context.Background(), callEnv.Request)

	assert.Equal(t, service.ErrUserNameIsEmpty, err)
}

func TestService_Register_DuplicateName_ReturnErrDuplicateName(t *testing.T) {
	serviceEnv := makeServiceEnv()
	serviceEnv.storage.
		On("GetByName", mock.Anything).
		Return(service.User{}, nil)

	callEnv := setupServcieRegisterEnv(serviceEnv)

	err := serviceEnv.Service.Register(context.Background(), callEnv.Request)

	assert.Equal(t, service.ErrDuplicateName, err)
}

func TestService_Register_UserStorage_GetByNameReturnErr_ReturnErr(t *testing.T) {
	serviceEnv := makeServiceEnv()
	serviceEnv.storage.
		On("GetByName", mock.Anything).
		Return(service.User{}, fmt.Errorf("some error"))

	callEnv := setupServcieRegisterEnv(serviceEnv)

	err := serviceEnv.Service.Register(context.Background(), callEnv.Request)

	assert.Error(t, err)
}

func TestService_Register_PasswordHasher_HashReturnErr_ReturnErr(t *testing.T) {
	serviceEnv := makeServiceEnv()
	serviceEnv.pwdHasher.
		On("Hash", mock.Anything).
		Return("", fmt.Errorf("some error"))

	callEnv := setupServcieRegisterEnv(serviceEnv)

	err := serviceEnv.Service.Register(context.Background(), callEnv.Request)

	assert.Error(t, err)
}

func TestService_Register_UserStorage_AddReturnErr_ReturnErr(t *testing.T) {
	serviceEnv := makeServiceEnv()
	serviceEnv.storage.
		On("Add", mock.Anything).
		Return(fmt.Errorf("some error"))

	callEnv := setupServcieRegisterEnv(serviceEnv)

	err := serviceEnv.Service.Register(context.Background(), callEnv.Request)

	assert.Error(t, err)
}

func TestService_SignIn_UserStorage_GetByNameReturnErr_ReturnErr(t *testing.T) {
	serviceEnv := makeServiceEnv()
	serviceEnv.storage.
		On("GetByName", mock.Anything).
		Return(service.User{}, fmt.Errorf("some error"))

	callEnv := setupServcieSignInEnv(serviceEnv)

	_, err := serviceEnv.Service.SignIn(context.Background(), callEnv.Request)

	assert.Error(t, err)
}

func TestService_SignIn_UserStorage_GetByNameReturnNotFoundErr_ReturnErrUserNameOrPasswordIncorrect(t *testing.T) {
	serviceEnv := makeServiceEnv()
	serviceEnv.storage.
		On("GetByName", mock.Anything).
		Return(service.User{}, apierror.EntityNotFoundErr)

	callEnv := setupServcieSignInEnv(serviceEnv)

	_, err := serviceEnv.Service.SignIn(context.Background(), callEnv.Request)

	assert.Equal(t, service.ErrUserNameOrPasswordIncorrect, err)
}

func TestService_SignIn_PwdHasher_VerifyReturnFalse_ReturnErrUserNameOrPasswordIncorrect(t *testing.T) {
	serviceEnv := makeServiceEnv()
	serviceEnv.pwdHasher.
		On("Verify", mock.Anything, mock.Anything).
		Return(false)

	callEnv := setupServcieSignInEnv(serviceEnv)

	_, err := serviceEnv.Service.SignIn(context.Background(), callEnv.Request)

	assert.Equal(t, service.ErrUserNameOrPasswordIncorrect, err)
}

func TestService_SignIn_TokenGenerator_MakeReturnErr_ReturnErr(t *testing.T) {
	serviceEnv := makeServiceEnv()
	serviceEnv.tokenG.
		On("Make", mock.Anything).
		Return("", fmt.Errorf("some error"))

	callEnv := setupServcieSignInEnv(serviceEnv)

	_, err := serviceEnv.Service.SignIn(context.Background(), callEnv.Request)

	assert.Error(t, err)
}

func TestService_SignIn_Seccess_ReturnVal(t *testing.T) {
	serviceEnv := makeServiceEnv()
	callEnv := setupServcieSignInEnv(serviceEnv)
	expected := service.SignInResp{
		Refresh:   "",
		Token:     callEnv.Token,
		TokenType: "vchat",
	}

	resp, err := serviceEnv.Service.SignIn(context.Background(), callEnv.Request)

	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestService_Join_Seccess_RetrunNil(t *testing.T) {
	serviceEnv := makeServiceEnv()
	callEnv := setupServcieJoinEnv(serviceEnv)

	err := serviceEnv.Service.Join(context.Background(), callEnv.Request)

	assert.NoError(t, err)
}

func TestService_Join_Seccess_UserUpdated(t *testing.T) {
	serviceEnv := makeServiceEnv()
	callEnv := setupServcieJoinEnv(serviceEnv)

	serviceEnv.Service.Join(context.Background(), callEnv.Request)

	serviceEnv.storage.AssertCalled(t, "Update", callEnv.UserAfter)
}

func TestService_Join_TokenInvalid_ReturnUnauthorizedRequestErr(t *testing.T) {
	serviceEnv := makeServiceEnv()
	serviceEnv.tokenG.
		On("Verify", mock.Anything).
		Return(service.Token{}, apierror.UnauthorizedRequestErr)

	callEnv := setupServcieJoinEnv(serviceEnv)

	err := serviceEnv.Service.Join(context.Background(), callEnv.Request)

	assert.Equal(t, apierror.UnauthorizedRequestErr, err)
}

func TestService_Join_StorageGetReturnErr_ReturnErr(t *testing.T) {
	serviceEnv := makeServiceEnv()
	serviceEnv.storage.
		On("Get", mock.Anything).
		Return(service.User{}, fmt.Errorf("some error"))

	callEnv := setupServcieJoinEnv(serviceEnv)

	err := serviceEnv.Service.Join(context.Background(), callEnv.Request)

	assert.Error(t, err)
}

func TestService_Join_GroupIDDuplicated_ReturnErrUnionIDDuplicated(t *testing.T) {
	serviceEnv := makeServiceEnv()

	callEnv := setupServcieJoinEnv(serviceEnv)
	callEnv.UserBefore.GroupIDs = append(callEnv.UserBefore.GroupIDs, callEnv.Request.UnionID)

	err := serviceEnv.Service.Join(context.Background(), callEnv.Request)

	assert.Equal(t, service.ErrUnionIDDuplicated, err)
}

func TestService_Join_StorageUpdateReturnErr_ReturnErr(t *testing.T) {
	serviceEnv := makeServiceEnv()
	serviceEnv.storage.
		On("Update", mock.Anything).
		Return(fmt.Errorf("some error"))

	callEnv := setupServcieJoinEnv(serviceEnv)

	err := serviceEnv.Service.Join(context.Background(), callEnv.Request)

	assert.Error(t, err)
}
