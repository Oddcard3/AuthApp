package messages

import (
	"time"
)

// Message chat message
type Message struct {
	ID      string
	Text    string
	Created time.Time
	ChatID  int
	Creator int
}
