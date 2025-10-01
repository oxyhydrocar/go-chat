package models

import (
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}

type MessageStore struct {
	messages []Message
}

func (ms *MessageStore) Add(msg Message) {
	ms.messages = append(ms.messages, msg)
}

func (ms *MessageStore) GetRecent(limit int) []Message {
	if len(ms.messages) <= limit {
		return ms.messages
	}
	return ms.messages[len(ms.messages)-limit:]
}
