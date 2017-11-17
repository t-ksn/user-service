package service

import (
	"net/http"

	"github.com/t-ksn/core-kit/apierror"
)

var (
	ErrDuplicateName = apierror.APIError{
		Code:       100,
		Message:    "User name already exist",
		StatusCode: http.StatusBadRequest,
	}

	ErrUserNameOrPasswordIncorrect = apierror.APIError{
		Code:       101,
		Message:    "User name or password incorrect",
		StatusCode: http.StatusBadRequest,
	}
	ErrPasswordLessThen4Chars = apierror.APIError{
		Code:       102,
		Message:    "Minimum password length 4",
		StatusCode: http.StatusBadRequest,
	}
	ErrUserNameIsEmpty = apierror.APIError{
		Code:       103,
		Message:    "User name is empty",
		StatusCode: http.StatusBadRequest,
	}
	ErrUnionIDDuplicated = apierror.APIError{
		Code:       104,
		Message:    "User already joined to union id",
		StatusCode: http.StatusBadRequest,
	}
)
