package bootstrap

import (
	"messages_handler/internal/config"
	"messages_handler/internal/messages_handler/domain/repository"
)

type RepositoryFactory struct {
	MessageRepository repository.IMessageRepository
}

func NewRepositoryFactory(cfg *config.Config) *RepositoryFactory {
	return &RepositoryFactory{
		MessageRepository: repository.NewMessageRepositoryImpl(cfg),
	}
}
