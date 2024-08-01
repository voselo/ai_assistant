package entity

// import "time"

type Message struct {
	ChannelID    string
	RefMessageId string
	CrmUserId    string
	CrmMessageId string
	ChatId       string
	ChatType     string
	Text         string
	// Timestamp time.Time
}
