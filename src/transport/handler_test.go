package transport_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/t-ksn/core-kit/apierror"
	"github.com/t-ksn/user-service/src/service"
	"github.com/t-ksn/user-service/src/transport"
)

type fakeService struct {
	mock.Mock
}

func (f *fakeService) Register(ctx context.Context, req service.CreateUserReq) error {
	args := f.Called(req)
	return args.Error(0)
}

func (f *fakeService) SignIn(ctx context.Context, req service.SignInReq) (service.SignInResp, error) {
	args := f.Called(req)
	return args.Get(0).(service.SignInResp), args.Error(1)
}

type BrokenBodyReader struct{}

func (r BrokenBodyReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("reader returned err")
}

func TestTable_Handler_ReadBody_ReturnErr_ReturnJSONInvalidErr(t *testing.T) {
	table := []struct {
		name string
		exec func(req *http.Request) error
	}{
		{
			name: "register",
			exec: func(req *http.Request) error {
				handler := transport.Make(nil)
				_, err := handler.Register(req)
				return err
			},
		},
		{
			name: "signIn",
			exec: func(req *http.Request) error {
				handler := transport.Make(nil)
				_, err := handler.SignIn(req)
				return err
			},
		},
	}

	for _, test := range table {
		t.Run(test.name, func(subt *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/read_body", BrokenBodyReader{})
			err := test.exec(req)
			assert.Equal(subt, apierror.JSONInvalidErr, err)
		})
	}
}

func TestTable_Handler_UnmarshalJSON_ReturnErr_ReturnJSONInvalidErr(t *testing.T) {
	table := []struct {
		name string
		exec func(req *http.Request) error
	}{
		{
			name: "register",
			exec: func(req *http.Request) error {
				handler := transport.Make(nil)
				_, err := handler.Register(req)
				return err
			},
		},
		{
			name: "signIn",
			exec: func(req *http.Request) error {
				handler := transport.Make(nil)
				_, err := handler.SignIn(req)
				return err
			},
		},
	}

	for _, test := range table {
		t.Run(test.name, func(subt *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/read_body", strings.NewReader("incorrect JSON 'oops "))
			err := test.exec(req)
			assert.Equal(subt, apierror.JSONInvalidErr, err)
		})
	}
}

func TestHandler_Register_ServiceExec_ReturnResult(t *testing.T) {
	expectedReq := service.CreateUserReq{
		Name:     "new user name",
		Password: "123",
	}
	reqSource := `{"name":"new user name","password":"123"}`
	expectedErr := fmt.Errorf("Expected error")

	fService := &fakeService{}
	fService.
		On("Register", expectedReq).
		Return(expectedErr)

	req, _ := http.NewRequest(http.MethodGet, "/read_body", strings.NewReader(reqSource))
	handler := transport.Make(fService)

	obj, err := handler.Register(req)

	assert.Nil(t, obj)
	assert.Equal(t, expectedErr, err)
}

func TestHandler_SignIn_ServiceExec_ReturnResult(t *testing.T) {
	expectedReq := service.SignInReq{
		Name:     "new user name",
		Password: "123",
	}
	reqSource := `{"name":"new user name","password":"123"}`
	expectedErr := fmt.Errorf("Expected error")
	expectedResp := service.SignInResp{
		Token:     "access_token",
		Refresh:   "refresh-tocken",
		TokenType: "token_type",
		ExpiredIn: 3600,
	}

	fService := &fakeService{}
	fService.
		On("SignIn", expectedReq).
		Return(expectedResp, expectedErr)

	req, _ := http.NewRequest(http.MethodGet, "/read_body", strings.NewReader(reqSource))
	handler := transport.Make(fService)

	obj, err := handler.SignIn(req)

	assert.Equal(t, expectedResp, obj)
	assert.Equal(t, expectedErr, err)
}
