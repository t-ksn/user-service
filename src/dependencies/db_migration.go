package dependencies

import (
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

var migrations = []func(db *mgo.Session) error{
	create_user_index_by_name,
}

func applyMigrations(db *mgo.Session) error {
	for _, m := range migrations {
		err := m(db)
		if err != nil {
			return err
		}
	}
	return nil
}

func create_user_index_by_name(db *mgo.Session) error {
	err := db.DB("").C("users").EnsureIndexKey("name")
	return errors.Wrapf(err, "Create index by name for users collection")
}
