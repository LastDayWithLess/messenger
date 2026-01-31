package repository

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/model"

	"gorm.io/gorm"
)

type MessengerRepository struct {
	db *gorm.DB
}

func NewMessengerRepository(db *gorm.DB) *MessengerRepository {
	return &MessengerRepository{db: db}
}

func (r *MessengerRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	return tx, nil
}

func (r *MessengerRepository) Commit(tx *gorm.DB, ctx context.Context) error {
	if err := tx.WithContext(ctx).Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *MessengerRepository) Rollback(tx *gorm.DB, ctx context.Context) error {
	if err := tx.WithContext(ctx).Rollback().Error; err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}

func (r *MessengerRepository) CreateChat(title string, tx *gorm.DB, ctx context.Context) (*model.Chats, error) {
	chat := &model.Chats{
		Title: title,
	}

	err := gorm.G[model.Chats](tx).Create(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	return chat, nil
}

func (r *MessengerRepository) CreateMessage(chatID uint, text string, tx *gorm.DB, ctx context.Context) (*model.Message, error) {
	message := &model.Message{
		Chat_ID: chatID,
		Text:    text,
	}

	err := gorm.G[model.Message](tx).Create(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return message, nil
}

func (r *MessengerRepository) DeleteChat(chatID uint, tx *gorm.DB, ctx context.Context) error {
	var chat model.Chats
	err := tx.WithContext(ctx).First(&chat, chatID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("chat with id %d not found: %w", chatID, err)
		}
		return fmt.Errorf("failed to find chat: %w", err)
	}

	result := tx.WithContext(ctx).Delete(&chat)
	if result.Error != nil {
		return fmt.Errorf("failed to delete chat: %w", result.Error)
	}

	return nil
}

func (r *MessengerRepository) GetChat(chatID uint, limit int, tx *gorm.DB, ctx context.Context) ([]model.MessageTimeDTO, error) {
	var chat model.Chats
	err := tx.WithContext(ctx).First(&chat, chatID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("chat with id %d not found: %w", chatID, err)
		}
		return nil, fmt.Errorf("failed to find chat: %w", err)
	}

	var messages []model.MessageTimeDTO
	result := tx.WithContext(ctx).
		Model(&model.Message{}).
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get chat messages: %w", result.Error)
	}

	for i := 0; i < len(messages)/2; i++ {
		messages[i], messages[len(messages)-i-1] = messages[len(messages)-i-1], messages[i]
	}

	return messages, nil
}
