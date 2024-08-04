package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"messages_handler/internal/config"
	"messages_handler/internal/messages_handler/domain/entity"
	"messages_handler/pkg/logging"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Interface
type IMessageRepository interface {
	HandleMessage(message entity.MessageEntity) error
	FindByChannelID(channelID string) ([]entity.MessageEntity, error)
}

// Implementation
type MessageRepository struct {
	cfg                *config.Config
	messages           map[string][]entity.MessageEntity
	timers             map[string]*time.Timer
	timerDuration      time.Duration
	firstMessageFactor int
	mu                 sync.Mutex
}

func NewMessageRepositoryImpl(cfg *config.Config) *MessageRepository {
	return &MessageRepository{
		cfg:                cfg,
		messages:           make(map[string][]entity.MessageEntity),
		timers:             make(map[string]*time.Timer),
		timerDuration:      10 * time.Second,
		firstMessageFactor: 3,
	}
}

func (repository *MessageRepository) HandleMessage(message entity.MessageEntity) error {
	logger := logging.GetLogger("trace")

	repository.mu.Lock()
	defer repository.mu.Unlock()

	// Save handled message to map, sorted by channelId
	repository.messages[message.ChannelId] = append(repository.messages[message.ChannelId], message)

	if timer, exists := repository.timers[message.ChannelId]; exists {
		timer.Stop()
	}

	duration := repository.timerDuration
	if len(repository.messages[message.ChannelId]) == 1 {
		duration = time.Duration(repository.firstMessageFactor) * repository.timerDuration
	}

	// Timer setup & Grouping messages after timer countdown
	repository.timers[message.ChannelId] = time.AfterFunc(duration,
		func() {
			repository.mu.Lock()
			defer repository.mu.Unlock()

			messages := repository.messages[message.ChannelId]

			var texts []string
			for _, msg := range messages {
				texts = append(texts, msg.Text)
			}

			grouppedMessage := messages[len(messages)-1]
			grouppedMessage.Text = strings.Join(texts, " ")

			// Clean up after sending
			logger.Infof("Timer expired for channel \n%s\nPerforming cleanup.\n", grouppedMessage.Text)
			delete(repository.messages, message.ChannelId)
			delete(repository.timers, message.ChannelId)

			// Converting to json
			var requestData = entity.MessageRequest{
				Messages: []entity.MessageEntity{grouppedMessage},
			}
			jsonData, err := json.Marshal(requestData)
			if err != nil {
				logger.Errorf("Error marshalling JSON: %v\n", err)
				return
			}

			// Api request
			url := repository.cfg.Webhook
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			if err != nil {
				logger.Errorf("Error creating request: %v\n", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				logger.Errorf("Error sending request: %v\n", err)
				return
			}
			defer resp.Body.Close()

			// Response
			logger.Infof("Response status: %s\n", resp.Status)

		},
	)

	return nil
}

func (repository *MessageRepository) FindByChannelID(channelID string) ([]entity.MessageEntity, error) {
	messages, exists := repository.messages[channelID]
	if !exists {
		return nil, errors.New("no messages found for channel")
	}
	return messages, nil
}
