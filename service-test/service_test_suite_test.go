package service_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/t-ksn/core-kit/debugclient"
	"github.com/t-ksn/core-kit/httpclient"

	"testing"
)

var userServiceClient UserServiceClient

func TestMain(m *testing.T) {
	var client httpclient.HTTPClient = http.DefaultClient
	if os.Getenv("DEBUG") != "" {
		client = &debugclient.DebugClient{Client: client}
	}

	var serviceAddress = os.Getenv("SERVICE_ADDRESS")
	if serviceAddress == "" {
		serviceAddress = "http://localhost:8080"
	}

	userServiceClient = UserServiceClient{
		Client: httpclient.VChatClient{
			ServiceAddress: serviceAddress,
			Client:         client,
		},
	}
}

func TestServiceTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User service Suite")
}

func safeGenString() string {
	var b [32]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}

type UserServiceClient struct {
	Client httpclient.VChatClient
}

func (c *UserServiceClient) Register(username, password string) error {
	return c.Client.Send(context.Background(), http.MethodPost, "/account", map[string]string{
		"name":     username,
		"password": password,
	}, nil)
}

type AccessToken struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
	Type    string `json:"type"`
}

func (c *UserServiceClient) SignIn(username, password string) (AccessToken, error) {
	var token AccessToken
	err := c.Client.Send(context.Background(), http.MethodPost, "/login", map[string]string{
		"name":     username,
		"password": password,
	}, &token)
	return token, err
}
