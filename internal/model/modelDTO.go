package model

import (
	"time"
)

type Chats struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	CreatedAt time.Time `gorm:"type:TIMESTAMPTZ;not null;index;default:CURRENT_TIMESTAMP(0)" json:"created_at"`
}

type Message struct {
	ID        uint      `gorm:"primaryKey;" json:"id"`
	Chat_ID   uint      `gorm:"not null;index" json:"chat_id"`
	Text      string    `gorm:"not null" json:"text"`
	CreatedAt time.Time `gorm:"type:TIMESTAMPTZ;not null;index;default:CURRENT_TIMESTAMP(0)" json:"created_at"`
}

type ErrorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type ChatDTO struct {
	Title string `json:"title"`
}

type MessageDTO struct {
	Text string `json:"text"`
}

type MessageTimeDTO struct {
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"time"`
}
