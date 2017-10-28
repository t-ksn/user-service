package transport

func Make(service Service) *Handler {
	return &Handler{servcie: service}
}
