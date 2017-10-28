package dependencies

import (
	"context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/pkg/errors"
	"github.com/t-ksn/core-kit/apierror"
	"github.com/t-ksn/user-service/src/service"
)

type MongoDBUserStorage struct {
	col *mgo.Collection
}

func (s *MongoDBUserStorage) GetByName(ctx context.Context, name string) (service.User, error) {
	var u service.User
	err := s.col.Find(bson.M{"name": name}).One(&u)
	if err == mgo.ErrNotFound {
		return u, apierror.EntityNotFoundErr
	}
	return u, errors.Wrapf(err, "MongoDBUserStorage.GetByName(%v)", name)
}

func (s *MongoDBUserStorage) Add(ctx context.Context, user service.User) error {
	err := s.col.Insert(user)
	return errors.Wrapf(err, "MongoDBUserStorage.Add(%#v)", user)
}
