package service

type User struct {
	ID           string   `bson:"_id"`
	Name         string   `bson:"name"`
	PasswordHash string   `bson:"pwd_hash"`
	GroupIDs     []string `bson:"group_ids"`
}

type Token struct {
	UserID string
}
