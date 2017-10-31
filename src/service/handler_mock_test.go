package service_test

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/t-ksn/core-kit/apierror"
	"github.com/t-ksn/user-service/src/service"
)

type fakeUserStorage struct {
	mock.Mock
}

func (f *fakeUserStorage) GetByName(ctx context.Context, name string) (service.User, error) {
	args := f.Called(name)
	return args.Get(0).(service.User), args.Error(1)
}

func (f *fakeUserStorage) Add(ctx context.Context, user service.User) error {
	args := f.Called(user)
	return args.Error(0)
}

type fakePasswordHasher struct {
	mock.Mock
}

func (f *fakePasswordHasher) Hash(password string) (string, error) {
	args := f.Called(password)
	return args.String(0), args.Error(1)
}

func (f *fakePasswordHasher) Verify(password, hash string) bool {
	args := f.Called(password, hash)
	return args.Bool(0)
}

type fakeTokenGenerator struct {
	mock.Mock
}

func (f *fakeTokenGenerator) Make(t service.Token) (string, error) {
	args := f.Called(t)
	return args.Get(0).(string), args.Error(1)
}

func (f *fakeTokenGenerator) Verify(t string) (service.Token, error) {
	args := f.Called(t)
	return args.Get(0).(service.Token), args.Error(1)
}

type serviceEnv struct {
	storage    *fakeUserStorage
	pwdHasher  *fakePasswordHasher
	tokenG     *fakeTokenGenerator
	Service    *service.Service
	idGenerate func() string
}

func makeServiceEnv() serviceEnv {
	result := serviceEnv{
		storage:    &fakeUserStorage{},
		pwdHasher:  &fakePasswordHasher{},
		tokenG:     &fakeTokenGenerator{},
		idGenerate: func() string { return "new_id" },
	}
	result.Service = service.Make(result.storage, result.pwdHasher, result.tokenG, result.idGenerate)

	return result
}

type serviceRegisterCallEnv struct {
	Request service.CreateUserReq
	User    service.User
}

func setupServcieRegisterEnv(se serviceEnv) serviceRegisterCallEnv {
	result := serviceRegisterCallEnv{
		Request: service.CreateUserReq{
			Name:     "new user name",
			Password: "1234",
		},
		User: service.User{
			ID:           se.idGenerate(),
			Name:         "new user name",
			PasswordHash: "pwd_hash",
		},
	}
	se.storage.
		On("GetByName", result.User.Name).
		Return(service.User{}, apierror.EntityNotFoundErr)

	se.pwdHasher.
		On("Hash", result.Request.Password).
		Return(result.User.PasswordHash, nil)
	se.storage.
		On("Add", result.User).
		Return(nil)

	return result
}

type serviceSignInEnv struct {
	Request service.SignInReq
	Token   string
	User    service.User
}

func setupServcieSignInEnv(se serviceEnv) serviceSignInEnv {
	result := serviceSignInEnv{
		Request: service.SignInReq{
			Name:     "user name",
			Password: "1234",
		},
		Token: "access_token",
		User: service.User{
			ID:           se.idGenerate(),
			Name:         "user name",
			PasswordHash: "pwd_hash",
		},
	}
	se.storage.
		On("GetByName", result.Request.Name).
		Return(result.User, nil)
	se.pwdHasher.
		On("Verify", result.Request.Password, result.User.PasswordHash).
		Return(true)
	se.tokenG.
		On("Make", service.Token{UserID: result.User.ID}).
		Return(result.Token, nil)

	return result
}
