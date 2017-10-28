package transport

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/t-ksn/core-kit/apierror"
	"github.com/t-ksn/user-service/src/service"
)

type Service interface {
	Register(ctx context.Context, req service.CreateUserReq) error
	SignIn(ctx context.Context, req service.SignInReq) (service.SignInResp, error)
}

type Handler struct {
	servcie Service
}

func (h *Handler) Register(req *http.Request) (interface{}, error) {
	var obj service.CreateUserReq
	err := unmarshal(req.Body, &obj)
	if err != nil {
		return nil, err
	}
	err = h.servcie.Register(req.Context(), obj)
	return nil, err
}

func (h *Handler) SignIn(req *http.Request) (interface{}, error) {
	var obj service.SignInReq
	err := unmarshal(req.Body, &obj)
	if err != nil {
		return nil, err
	}
	return h.servcie.SignIn(req.Context(), obj)
}

func unmarshal(reader io.Reader, obj interface{}) error {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return apierror.JSONInvalidErr
	}
	err = json.Unmarshal(body, obj)
	if err != nil {
		return apierror.JSONInvalidErr
	}
	return nil
}
