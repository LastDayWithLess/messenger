package service

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/loggerzap"
	"messenger/internal/model"
	"messenger/internal/repository"

	"go.uber.org/zap"
)

var (
	ErrLen = errors.New("incorrect length title")
)

type Logger interface {
	Error(err string, field ...zap.Field)
	Info(msg string, field ...zap.Field)
}

type MessengerServiceInt interface {
	CreateChat(title string, ctx context.Context) (*model.Chats, error)
	CreateMessage(text string, chatID uint, ctx context.Context) (*model.Message, error)
	GetChat(chatID uint, limit int, ctx context.Context) (map[uint][]model.MessageTimeDTO, error)
	DeleteChat(chatID uint, ctx context.Context) error
	GetLogger() Logger
}

type MessengerService struct {
	db     *repository.MessengerRepository
	logger Logger
}

func NewMessengerService(db *repository.MessengerRepository) (*MessengerService, error) {
	logger, err := loggerzap.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &MessengerService{
		db:     db,
		logger: logger,
	}, nil
}

func (s *MessengerService) CreateChat(title string, ctx context.Context) (*model.Chats, error) {
	if len(title) == 0 || len(title) > 200 {
		err := fmt.Errorf("incorrect length title: %d", len(title))
		s.logger.Error("Validation failed", loggerzap.ErrorField(err))
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transaction", loggerzap.ErrorField(err))
		return nil, err
	}

	defer func() {
		if err != nil {
			if err := s.db.Rollback(tx, ctx); err != nil {
				s.logger.Error("Failed to rollback transaction", loggerzap.ErrorField(err))
			}
		}
	}()

	chat, err := s.db.CreateChat(title, tx, ctx)
	if err != nil {
		return nil, err
	}

	if err = s.db.Commit(tx, ctx); err != nil {
		s.logger.Error("Failed to commit transaction", loggerzap.ErrorField(err))
		return nil, err
	}

	return chat, nil
}

func (s *MessengerService) CreateMessage(text string, chatID uint, ctx context.Context) (*model.Message, error) {
	if len(text) == 0 || len(text) > 5000 {
		s.logger.Error("Validation failed", loggerzap.ErrorField(ErrLen))

		return nil, ErrLen
	}

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Fieled to begin transaction", loggerzap.ErrorField(err))

		return nil, err
	}

	defer func() {
		if err != nil {
			if err := s.db.Rollback(tx, ctx); err != nil {
				s.logger.Error("Failed to rollback transaction", loggerzap.ErrorField(err))
			}
		}
	}()

	message, err := s.db.CreateMessage(chatID, text, tx, ctx)
	if err != nil {
		return nil, err
	}

	if err = s.db.Commit(tx, ctx); err != nil {
		s.logger.Error("Failed to commit transaction", loggerzap.ErrorField(err))
		return nil, err
	}

	return message, nil
}

func (s *MessengerService) DeleteChat(chatID uint, ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Fieled to begin transaction", loggerzap.ErrorField(err))

		return err
	}

	defer func() {
		if err != nil {
			if err = s.db.Rollback(tx, ctx); err != nil {
				s.logger.Error("Failed to rollback transaction", loggerzap.ErrorField(err))
			}
		}
	}()

	if err = s.db.DeleteChat(chatID, tx, ctx); err != nil {
		return err
	}

	if err = s.db.Commit(tx, ctx); err != nil {
		s.logger.Error("Failed to commit transaction", loggerzap.ErrorField(err))
		return err
	}

	return nil
}

func (s *MessengerService) GetChat(chatID uint, limit int, ctx context.Context) (map[uint][]model.MessageTimeDTO, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Fieled to begin transaction", loggerzap.ErrorField(err))

		return nil, err
	}

	defer func() {
		if err != nil {
			if err = s.db.Rollback(tx, ctx); err != nil {
				s.logger.Error("Failed to rollback transaction", loggerzap.ErrorField(err))
			}
		}
	}()

	messages, err := s.db.GetChat(chatID, limit, tx, ctx)
	if err != nil {
		return nil, err
	}

	if err = s.db.Commit(tx, ctx); err != nil {
		s.logger.Error("Failed to commit transaction", loggerzap.ErrorField(err))
		return nil, err
	}

	mapMessage := map[uint][]model.MessageTimeDTO{}
	mapMessage[chatID] = messages

	return mapMessage, nil
}

func (s *MessengerService) GetLogger() Logger {
	return s.logger
}
