package service

func Make(s UserStorage, ph PasswordHasher, tg TokenGenerator) *Service {
	return &Service{
		storage:   s,
		tokenG:    tg,
		pwdHasher: ph,
	}
}
