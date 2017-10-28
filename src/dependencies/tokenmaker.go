package dependencies

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/t-ksn/core-kit/apierror"
	"github.com/t-ksn/user-service/src/service"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

func init() {
	jwt.TimeFunc = func() time.Time {
		return time.Now().UTC()
	}
}

type TokenMaker struct {
	secret []byte
}

func (m *TokenMaker) Make(t service.Token) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  t.UserID,
		"exp": t.Exp,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(m.secret)
	return tokenString, errors.Wrap(err, "TokenMaker.Make")
}

func (m *TokenMaker) Verify(tokenString string) (service.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return m.secret, nil
	})
	if err != nil {
		return service.Token{}, errors.Wrap(err, "TokenMaker.Verify")
	}

	if token.Valid {
		return service.Token{}, apierror.UnauthorizedRequestErr
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		t := service.Token{}
		t.UserID, _ = claims["id"].(string)
		t.Exp, _ = claims["exp"].(int64)
		return t, nil
	}

	return service.Token{}, fmt.Errorf("TokenMaker.Verify: unsupported claims format %#v", token.Claims)
}
