package model

import (
	"time"
)

type ChatDetails struct {
	Messages []MessageModel
	Timer    *time.Timer
}
