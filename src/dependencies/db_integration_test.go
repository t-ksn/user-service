// +build integration

package dependencies_test

import (
	"context"
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/t-ksn/core-kit/apierror"
	"github.com/t-ksn/user-service/src/service"
)

func TestUserStorage_Add(t *testing.T) {
	err := userStorage.Add(context.Background(), service.User{
		Name:         "test_add_user_storage",
		ID:           uuid.NewV4().String(),
		PasswordHash: "1234",
	})

	assert.NoError(t, err)
}

func TestUserStorage_GetByName(t *testing.T) {
	expected := service.User{
		Name:         uuid.NewV4().String(),
		ID:           uuid.NewV4().String(),
		PasswordHash: "1234",
		UnionIDs:     []string{},
	}
	err := userStorage.Add(context.Background(), expected)
	assert.NoError(t, err)

	actual, err := userStorage.GetByName(context.Background(), expected.Name)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUserStorage_GetByName_UserNotFound_ReturnNotFoundErr(t *testing.T) {
	_, err := userStorage.GetByName(context.Background(), uuid.NewV4().String())
	assert.Equal(t, apierror.EntityNotFoundErr, err)
}

func TestUserStorage_Update(t *testing.T) {
	user := service.User{
		Name:         "test_update_user_storage",
		ID:           uuid.NewV4().String(),
		PasswordHash: "1234",
		UnionIDs:     []string{},
	}

	err := userStorage.Add(context.Background(), user)
	assert.NoError(t, err)

	user.Name = "new name"
	err = userStorage.Update(context.Background(), user)
	assert.NoError(t, err)

	actual, err := userStorage.Get(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user, actual)
}

func TestUserStorage_Get(t *testing.T) {
	expected := service.User{
		Name:         "test_get_user_storage",
		ID:           uuid.NewV4().String(),
		PasswordHash: "1234",
		UnionIDs:     []string{},
	}
	err := userStorage.Add(context.Background(), expected)
	assert.NoError(t, err)

	actual, err := userStorage.Get(context.Background(), expected.ID)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
