package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/t-ksn/core-kit/apierror"
)

type UserStorage interface {
	GetByName(ctx context.Context, name string) (User, error)
	Get(ctx context.Context, id string) (User, error)

	Add(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
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
	storage    UserStorage
	pwdHasher  PasswordHasher
	tokenG     TokenGenerator
	idGenerate func() string
}

func (s *Service) Register(ctx context.Context, req CreateUserReq) error {
	if len(req.Password) < 4 {
		return ErrPasswordLessThen4Chars
	}
	if req.Name == "" {
		return ErrUserNameIsEmpty
	}

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

	u := User{Name: req.Name, PasswordHash: pwdHash, ID: s.idGenerate()}
	err = s.storage.Add(ctx, u)
	return errors.WithMessage(err, "Servcie.Register")
}

func (s *Service) SignIn(ctx context.Context, req SignInReq) (SignInResp, error) {
	u, err := s.storage.GetByName(ctx, req.Name)
	if err == apierror.EntityNotFoundErr {
		return SignInResp{}, ErrUserNameOrPasswordIncorrect
	}
	if err != nil {
		return SignInResp{}, errors.WithMessage(err, "Service.SignIn")
	}
	verified := s.pwdHasher.Verify(req.Password, u.PasswordHash)
	if !verified {
		return SignInResp{}, ErrUserNameOrPasswordIncorrect
	}
	token, err := s.tokenG.Make(Token{UserID: u.ID})
	if err != nil {
		return SignInResp{}, errors.WithMessage(err, "Service.SignIn")
	}
	return SignInResp{
		Token:     token,
		TokenType: "vchat",
	}, nil
}

func (s *Service) Join(ctx context.Context, req Join2Req) error {
	token, err := s.tokenG.Verify(req.Token)
	if err != nil {
		return apierror.UnauthorizedRequestErr
	}

	user, err := s.storage.Get(ctx, token.UserID)
	if err != nil {
		return errors.WithMessage(err, "Service.Join2Group")
	}

	for i := 0; i < len(user.GroupIDs); i++ {
		if user.GroupIDs[i] == req.UnionID {
			return ErrUnionIDDuplicated
		}
	}

	user.GroupIDs = append(user.GroupIDs, req.UnionID)
	err = s.storage.Update(ctx, user)

	return errors.WithMessage(err, "Service.Join2Group")
}
