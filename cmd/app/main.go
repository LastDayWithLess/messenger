package main

import (
	"messenger/config"
	"messenger/internal/repository"
	"messenger/internal/service"
	"messenger/internal/transport"
)

func main() {
	err := config.Init(".env")
	if err != nil {
		return
	}

	v, err := config.LoadConfigDB()
	if err != nil {
		return
	}

	con, err := repository.NewConnection(v)
	if err != nil {
		return
	}

	test := repository.NewMessengerRepository(con)
	testServ, err := service.NewMessengerService(test)
	if err != nil {
		return
	}

	h := transport.NewHandlerMessenger(testServ)
	s := transport.NewServerMessenger(h)

	err = s.StartServer()
	if err != nil {
		return
	}
}
