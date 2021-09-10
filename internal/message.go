package internal

import "time"

type Message struct {
	Timestamp    time.Time `json:"timestamp"`
	Message      string    `json:"message"`
	Room         string    `json:"room"`
	DisplayName  string    `json:"displayName"`
	Type         string    `json:"type"`
	Participants int       `json:"participants"`
}
