package entity

import "time"

// import "time"

type MessageEntity struct {
	MessageId string `json:"messageId"`
	DateTime time.Time `json:"dateTime"`
	ChannelId string `json:"channelId"`
	ChannelType string `json:"chatType"`
	ChatId string `json:"chatId"`
	Type string `json:"type"`
	IsEcho bool `json:"isEcho"`
	Conact ContactEntity `json:"contact"`
	Text string `json:"text"`
	Status string `json:"status"`

}

