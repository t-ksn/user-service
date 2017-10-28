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
)
