package service

func Make(s UserStorage, ph PasswordHasher, tg TokenGenerator, idGenerate func() string) *Service {
	return &Service{
		storage:    s,
		tokenG:     tg,
		pwdHasher:  ph,
		idGenerate: idGenerate,
	}
}
