package model

import "time"


type MessageModel struct {
	MessageId   string        `json:"messageId"`
	DateTime    time.Time     `json:"dateTime"`
	ChannelId   string        `json:"channelId"`
	ChannelType string        `json:"chatType"`
	ChatId      string        `json:"chatId"`
	Type        string        `json:"type"`
	IsEcho      bool          `json:"isEcho"`
	Conact      ContactModel `json:"contact"`
	Text        string        `json:"text"`
	Status      string        `json:"status"`
}

type ContactModel struct {
	Name      string
	AvatarUri string
	Username  string
}
