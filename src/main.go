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
	tm := dependencies.MakeTokenMaker(secreatKey)
	bs := service.Make(us, ph, tm)
	h := transport.Make(bs)

	cs := core.NewService(core.Name("user-service"),
		core.Port(port),
		core.Version(version))

	cs.Post("/user", h.Register)
	cs.Post("/login", h.SignIn)

	cs.Run()
}
