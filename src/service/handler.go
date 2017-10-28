package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/t-ksn/core-kit/apierror"
)

type UserStorage interface {
	GetByName(ctx context.Context, name string) (User, error)
	Add(ctx context.Context, user User) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}

type TokenGenerator interface {
	Make(t Token) (string, error)
	Verify(string) (Token, error)
}

type Service struct {
	storage   UserStorage
	pwdHasher PasswordHasher
	tokenG    TokenGenerator
}

func (s *Service) Register(ctx context.Context, req CreateUserReq) error {
	_, err := s.storage.GetByName(ctx, req.Name)
	if err == nil {
		return ErrDuplicateName
	}

	if err != apierror.EntityNotFoundErr {
		return errors.WithMessage(err, "Service.Register")
	}

	pwdHash, err := s.pwdHasher.Hash(req.Password)
	if err != nil {
		return errors.WithMessage(err, "Servcie.Register")
	}

	u := User{Name: req.Name, PasswordHash: pwdHash, ID: uuid.NewV4().String()}
	err = s.storage.Add(ctx, u)
	return errors.WithMessage(err, "Servcie.Register")
}

const tokenTTL = time.Hour

func (s *Service) SignIn(ctx context.Context, req SignInReq) (SignInResp, error) {
	u, err := s.storage.GetByName(ctx, req.Name)
	if err == apierror.EntityNotFoundErr {
		return SignInResp{}, ErrUserNameOrPasswordIncorrect
	}
	if err != nil {
		return SignInResp{}, errors.WithMessage(err, "Service.SignIn")
	}
	exp := time.Now().UTC().Add(tokenTTL).Unix()
	token, err := s.tokenG.Make(Token{UserID: u.ID, Exp: exp})
	if err != nil {
		return SignInResp{}, errors.WithMessage(err, "Service.SignIn")
	}
	return SignInResp{
		Token:     token,
		ExpiredIn: exp,
		TokenType: "vchat",
	}, nil

}
