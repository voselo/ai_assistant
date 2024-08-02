package service

import (
	"sync"
	"time"

	"messages_handler/internal/domain/entity"
	"messages_handler/internal/domain/repository"
)

type MessageService struct {
	repo         repository.MessageRepository
	timeWindow   time.Duration
	messageCache map[string][]entity.MessageEntity
	mu           sync.Mutex
}

func NewMessageService(repo repository.MessageRepository) *MessageService {

	duration := time.Duration(5 * time.Second)

	return &MessageService{
		repo:         repo,
		timeWindow:   duration,
		messageCache: make(map[string][]entity.MessageEntity),
	}
}

func (s *MessageService) AddMessage(channelID, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// message := entity.Message{
	// 	ChannelID: channelID,
	// 	Content:   content,
	// 	Timestamp: time.Now(),
	// }

	// s.messageCache[channelID] = append(s.messageCache[channelID], message)

	// go s.processMessages(channelID)
	return nil
}

func (s *MessageService) processMessages(channelID string) {
	time.Sleep(s.timeWindow)

	s.mu.Lock()
	defer s.mu.Unlock()

	messages := s.messageCache[channelID]
	if len(messages) > 0 {

		// err := s.repo.SaveBatch(messages)
		// if err == nil {
		// 	delete(s.messageCache, channelID)
		// }
	}
}

func (service *MessageService) GetMessages(channelID string) ([]entity.MessageEntity, error) {
	return service.repo.FindByChannelID(channelID)
}
