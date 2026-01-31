package transport

import (
	"context"
	"encoding/json"
	"errors"
	"messenger/internal/model"
	"messenger/internal/service"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type HandlerMessenger struct {
	serv service.MessengerServiceInt
}

func NewHandlerMessenger(serv service.MessengerServiceInt) *HandlerMessenger {
	return &HandlerMessenger{
		serv: serv,
	}
}

func (h *HandlerMessenger) SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	errDTO := model.ErrorDTO{
		Message: message,
		Time:    time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(errDTO); err != nil {
		h.serv.GetLogger().Error(err.Error())
	}
}

func (h *HandlerMessenger) HendleCreateChat(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if !strings.Contains(contentType, "application/json") {
		h.SendErrorResponse(w, http.StatusUnsupportedMediaType, "Contetnt-Type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var chatDTO model.ChatDTO

	if err := json.NewDecoder(r.Body).Decode(&chatDTO); err != nil {
		h.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chat, err := h.serv.CreateChat(chatDTO.Title, ctx)
	if err != nil {
		h.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(chat); err != nil {
		h.serv.GetLogger().Error(err.Error())
	}
}

func (h *HandlerMessenger) HendleCreateMessage(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if !strings.Contains(contentType, "application/json") {
		h.SendErrorResponse(w, http.StatusUnsupportedMediaType, "Contetnt-Type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var messageDTO model.MessageDTO

	if err := json.NewDecoder(r.Body).Decode(&messageDTO); err != nil {
		h.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	chatID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(chatID)
	if err != nil {
		h.SendErrorResponse(w, http.StatusBadRequest, "id chat must be integer")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message, err := h.serv.CreateMessage(messageDTO.Text, uint(id), ctx)
	if err != nil {
		if !errors.Is(err, service.ErrLen) {
			h.SendErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			h.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	w.Header().Set("Content_Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(message); err != nil {
		h.serv.GetLogger().Error(err.Error())
	}
}

func (h *HandlerMessenger) HendleGetMessages(w http.ResponseWriter, r *http.Request) {
	chatID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(chatID)
	if err != nil {
		h.SendErrorResponse(w, http.StatusBadRequest, "id chat must be integer")
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil && limit != 0 {
		h.SendErrorResponse(w, http.StatusBadRequest, "id limit must be integer")
	}

	if limit == 0 {
		limit = 20
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	messages, err := h.serv.GetChat(uint(id), limit, ctx)
	if err != nil {
		h.SendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content_Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(messages); err != nil {
		h.serv.GetLogger().Error(err.Error())
	}
}

func (h *HandlerMessenger) HandleDeleteChat(w http.ResponseWriter, r *http.Request) {
	chatID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(chatID)
	if err != nil {
		h.SendErrorResponse(w, http.StatusBadRequest, "id chat must be integer")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.serv.DeleteChat(uint(id), ctx)

	w.Header().Set("Content_Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
