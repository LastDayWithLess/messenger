package transport

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerMessenger struct {
	httpHandler *HandlerMessenger
}

func NewServerMessenger(httpHandler *HandlerMessenger) *ServerMessenger {
	return &ServerMessenger{
		httpHandler: httpHandler,
	}
}

func (s *ServerMessenger) StartServer() error {
	router := mux.NewRouter()

	router.Path("/chats").HandlerFunc(s.httpHandler.HendleCreateChat).Methods("POST")
	router.Path("/chats/{id}/messages").HandlerFunc(s.httpHandler.HendleCreateMessage).Methods("POST")
	router.Path("/chats/{id}").HandlerFunc(s.httpHandler.HendleGetMessages).Methods("GET")
	router.Path("/chats/{id}").HandlerFunc(s.httpHandler.HandleDeleteChat).Methods("DELETE")

	if err := http.ListenAndServe(":8080", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
