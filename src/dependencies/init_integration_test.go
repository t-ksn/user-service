// +build integration

package dependencies_test

import (
	"os"
	"testing"

	"github.com/t-ksn/user-service/src/dependencies"
)

var userStorage *dependencies.MongoDBUserStorage

func TestMain(m *testing.M) {
	dbConnectionStr := os.Getenv("TEST_DB")
	if dbConnectionStr == "" {
		dbConnectionStr = "mongodb://localhost/test_user_service"
	}

	var err error
	userStorage, err = dependencies.MakeUserStorage(dbConnectionStr)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
