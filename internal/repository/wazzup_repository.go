package repository

import (
	model "ai_assistant/internal/model/messages"
	"ai_assistant/pkg/logging"
	"bytes"
	"encoding/json"
	"net/http"

	"strings"
	"sync"
	"time"
)

type WazzupRepository struct {
	chats              map[string]map[string]*model.ChatDetails
	timerDuration      time.Duration
	firstMessageFactor int
	mu                 sync.Mutex
}

func NewWazzupRepository() *WazzupRepository {
	return &WazzupRepository{
		chats:              make(map[string]map[string]*model.ChatDetails),
		timerDuration:      5 * time.Second,
		firstMessageFactor: 3,
	}
}

func (repository *WazzupRepository) ProcessMessage(uid string, message model.MessageModel, customerRepo *CustomersRepository) {
	logger := logging.GetLogger("Info")

	repository.mu.Lock()
	defer repository.mu.Unlock()

	// Sort message by channelId
	channelChats, ok := repository.chats[message.ChannelId]
	if !ok {
		channelChats = make(map[string]*model.ChatDetails)
		repository.chats[message.ChannelId] = channelChats
	}

	// Sort message by chatId
	chatDetails, ok := channelChats[message.ChatId]
	if !ok {
		chatDetails = &model.ChatDetails{
			Messages: make([]model.MessageModel, 0),
		}
		channelChats[message.ChatId] = chatDetails
	}

	// Save handled message to map, sorted by chatId
	chatDetails.Messages = append(chatDetails.Messages, message)

	if chatDetails.Timer != nil {
		chatDetails.Timer.Stop()
	}

	// First message duration factor
	duration := repository.timerDuration
	if len(chatDetails.Messages) == 1 {
		duration *= time.Duration(repository.firstMessageFactor)
	}

	// Timer after func
	chatDetails.Timer = time.AfterFunc(duration, func() {
		repository.mu.Lock()
		defer repository.mu.Unlock()

		// Groupping messages
		var texts []string
		for _, msg := range chatDetails.Messages {
			texts = append(texts, msg.Text)
		}

		grouppedMessage := chatDetails.Messages[len(chatDetails.Messages)-1]
		grouppedMessage.Text = strings.Join(texts, " ")

		logger.Infof("Message gropped: %s\n", grouppedMessage.Text)

		// Cleanup RAM
		delete(channelChats, message.ChatId)
		if len(channelChats) == 0 {
			delete(repository.chats, message.ChannelId)
		}

		logger.Info("Perfoming sending message")

		// Sending [Post] Request with groupped message to customer
		go func(grouppedMessage model.MessageModel) {

			// Get customer data by id
			customer, err := customerRepo.GetById(uid)
			if err != nil {
				logger.Error("Error creating HTTP request: ", err)
				return
			}

			// Serialise data
			jsonData, err := json.Marshal(grouppedMessage)
			if err != nil {
				logger.Error("Error marshalling JSON: ", err)
				return
			}

			url := customer.WazzupUri
			request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			if err != nil {
				logger.Error("Error creating HTTP request: ", err)
				return
			}

			request.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(request)
			if err != nil {
				logger.Error("Error sending HTTP request: ", err)
				return
			}
			defer resp.Body.Close()

			// Response processing
			if resp.StatusCode == http.StatusOK {
				logger.Info("Message sent successfully")
			} else {
				logger.Errorf("Failed to send message, status code: %d", resp.StatusCode)
			}
		}(grouppedMessage)
	})

}
