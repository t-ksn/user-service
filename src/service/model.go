package service

type User struct {
	ID           string `bson:"_id"`
	Name         string `bson:"name"`
	PasswordHash string `bson:"pwd_hash"`
}

type Token struct {
	UserID string
}
