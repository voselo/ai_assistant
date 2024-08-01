package repository

import (
	"errors"
	"messages_handler/internal/domain/entity"
)

// Interface
type MessageRepository interface {
	Save(message entity.Message) error
	FindByChannelID(channelID string) ([]entity.Message, error)
}

// Implementation
type InMemoryMessageRepository struct {
	messages map[string][]entity.Message
}

func NewInMemoryMessageRepository() *InMemoryMessageRepository {
	return &InMemoryMessageRepository{
		messages: make(map[string][]entity.Message),
	}
}

func (repository *InMemoryMessageRepository) Save(message entity.Message) error {
	repository.messages[message.ChannelID] = append(repository.messages[message.ChannelID], message)
	return nil
}

func (repository *InMemoryMessageRepository) FindByChannelID(channelID string) ([]entity.Message, error) {
	messages, exists := repository.messages[channelID]
	if !exists {
		return nil, errors.New("no messages found for channel")
	}
	return messages, nil
}
