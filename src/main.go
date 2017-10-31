package main

import (
	"log"

	"github.com/namsral/flag"

	"github.com/t-ksn/core-kit"
	"github.com/t-ksn/user-service/src/dependencies"
	"github.com/t-ksn/user-service/src/service"
	"github.com/t-ksn/user-service/src/transport"
)

var (
	version string
)
var (
	port            int
	dbConnectionStr string
	secreatKey      string
)

func init() {
	flag.IntVar(&port, "port", 80, "Port of listening")
	flag.StringVar(&dbConnectionStr, "db", "", "Database connection string")
	flag.StringVar(&secreatKey, "key", "", "Seacret key")
}

func main() {
	ph := dependencies.MakePasswordHasher()
	us, err := dependencies.MakeUserStorage(dbConnectionStr)
	if err != nil {
		log.Fatal(err)
	}
	tg := dependencies.MakeTokenGenerator(secreatKey)
	bs := service.Make(us, ph, tg)
	h := transport.Make(bs)

	cs := corekit.NewService(corekit.Name("user-service"),
		corekit.Port(port),
		corekit.Version(version))

	cs.Post("/user", h.Register)
	cs.Post("/login", h.SignIn)

	cs.Run()
}
