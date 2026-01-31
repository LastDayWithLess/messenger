package mocks

import (
	"context"
	"messenger/internal/model"
	"messenger/internal/service"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Error(err string, field ...zap.Field) {
	m.Called(err, field)
}

func (m *MockLogger) Info(msg string, field ...zap.Field) {
	m.Called(msg, field)
}

type MockMessengerService struct {
	mock.Mock
	logger *MockLogger
}

func NewMockMessengerService() *MockMessengerService {
	return &MockMessengerService{
		logger: &MockLogger{},
	}
}

func (m *MockMessengerService) CreateChat(title string, ctx context.Context) (*model.Chats, error) {
	args := m.Called(title, ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chats), args.Error(1)
}

func (m *MockMessengerService) CreateMessage(text string, chatID uint, ctx context.Context) (*model.Message, error) {
	args := m.Called(text, chatID, ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Message), args.Error(1)
}

func (m *MockMessengerService) GetChat(chatID uint, limit int, ctx context.Context) (map[uint][]model.MessageTimeDTO, error) {
	args := m.Called(chatID, limit, ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[uint][]model.MessageTimeDTO), args.Error(1)
}

func (m *MockMessengerService) DeleteChat(chatID uint, ctx context.Context) error {
	args := m.Called(chatID, ctx)
	return args.Error(0)
}

func (m *MockMessengerService) GetLogger() service.Logger {
	return m.logger
}
