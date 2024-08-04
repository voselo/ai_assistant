package entity

type MessageRequest struct {
	Messages []MessageEntity `json:"messages"`
}