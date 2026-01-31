package transport

import (
	"bytes"
	"encoding/json"
	"messenger/internal/model"
	"messenger/tests/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_HendleCreateChat_Success(t *testing.T) {
	mockService := new(mocks.MockMessengerService)
	handler := &HandlerMessenger{
		serv: mockService,
	}

	expectedChat := &model.Chats{
		ID:    1,
		Title: "Test Chat",
	}

	mockService.On("CreateChat", "Test Chat", mock.AnythingOfType("*context.timerCtx")).
		Return(expectedChat, nil)

	chatDTO := model.ChatDTO{Title: "Test Chat"}
	body, _ := json.Marshal(chatDTO)

	req := httptest.NewRequest("POST", "/chats", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HendleCreateChat(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var actualChat model.Chats
	err := json.Unmarshal(rr.Body.Bytes(), &actualChat)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), actualChat.ID)
	assert.Equal(t, "Test Chat", actualChat.Title)

	mockService.AssertExpectations(t)
}

func TestHandler_HendleCreateMessage_InvalidChatID(t *testing.T) {
	mockService := new(mocks.MockMessengerService)
	handler := &HandlerMessenger{
		serv: mockService,
	}

	messageDTO := model.MessageDTO{Text: "Hello World"}
	body, _ := json.Marshal(messageDTO)

	req := httptest.NewRequest("POST", "/chats/invalid/messages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	rr := httptest.NewRecorder()

	handler.HendleCreateMessage(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_HendleGetMessages_WithLimit(t *testing.T) {
	mockService := new(mocks.MockMessengerService)
	handler := &HandlerMessenger{
		serv: mockService,
	}

	expectedMessages := map[uint][]model.MessageTimeDTO{
		1: {
			{
				Text:      "Message 1",
				CreatedAt: time.Now(),
			},
			{
				Text:      "Message 2",
				CreatedAt: time.Now(),
			},
		},
	}

	mockService.On("GetChat", uint(1), 10, mock.AnythingOfType("*context.timerCtx")).
		Return(expectedMessages, nil)

	req := httptest.NewRequest("GET", "/chats/1/messages?limit=10", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handler.HendleGetMessages(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var messages map[uint][]model.MessageTimeDTO
	err := json.Unmarshal(rr.Body.Bytes(), &messages)
	assert.NoError(t, err)
	assert.Len(t, messages[1], 2)
	assert.Equal(t, "Message 1", messages[1][0].Text)
	assert.Equal(t, "Message 2", messages[1][1].Text)

	mockService.AssertExpectations(t)
}
