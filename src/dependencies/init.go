package dependencies

import "gopkg.in/mgo.v2"

func MakeUserStorage(dbConnectionStr string) (*MongoDBUserStorage, error) {
	s, err := mgo.Dial(dbConnectionStr)
	if err != nil {
		return nil, err
	}
	err = applayMigration(s)
	if err != nil {
		return nil, err
	}
	return &MongoDBUserStorage{col: s.DB("").C("users")}, nil
}

func MakePasswordHasher() *PasswordHasher {
	return &PasswordHasher{}
}

func MakeTokenMaker(seacretKey string) *TokenMaker {
	return &TokenMaker{secret: []byte(seacretKey)}
}
