package repository

import (
	"encoding/json"
	"messages_handler/internal/wazzup/model"
	"messages_handler/pkg/logging"

	"strings"
	"sync"
	"time"
)

type WazzupRepository struct {
	messages           map[string][]model.MessageModel
	timers             map[string]*time.Timer
	timerDuration      time.Duration
	firstMessageFactor int
	mu                 sync.Mutex
}

func New() *WazzupRepository {
	return &WazzupRepository{
		messages:           make(map[string][]model.MessageModel),
		timers:             make(map[string]*time.Timer),
		timerDuration:      10 * time.Second,
		firstMessageFactor: 3,
	}
}

func (repository *WazzupRepository) ProcessMessage(message model.MessageModel) {
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
			var requestData = model.MessageRequest{
				Messages: []model.MessageModel{grouppedMessage},
			}
			jsonData, err := json.Marshal(requestData)
			if err != nil {
				logger.Errorf("Error marshalling JSON: %v\n", err)
				return
			}

			logger.Info(jsonData)

		},
	)

}
