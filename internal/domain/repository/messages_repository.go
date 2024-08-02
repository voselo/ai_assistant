package repository

import (
	"errors"
	"messages_handler/internal/domain/entity"
)

// Interface
type MessageRepository interface {
	Save(message entity.MessageEntity) error
	FindByChannelID(channelID string) ([]entity.MessageEntity, error)
}

// Implementation
type InMemoryMessageRepository struct {
	messages map[string][]entity.MessageEntity
}

func NewInMemoryMessageRepository() *InMemoryMessageRepository {
	return &InMemoryMessageRepository{
		messages: make(map[string][]entity.MessageEntity),
	}
}

func (repository *InMemoryMessageRepository) Save(message entity.MessageEntity) error {
	repository.messages[message.ChannelID] = append(repository.messages[message.ChannelID], message)
	return nil
}

func (repository *InMemoryMessageRepository) FindByChannelID(channelID string) ([]entity.MessageEntity, error) {
	messages, exists := repository.messages[channelID]
	if !exists {
		return nil, errors.New("no messages found for channel")
	}
	return messages, nil
}
