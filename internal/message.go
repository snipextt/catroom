package internal

import "time"

type Message struct {
	Timestamp   time.Time `json:"timestamp,omitempty"`
	Message     string    `json:"message,omitempty"`
	Room        string    `json:"room,omitempty"`
	DisplayName string    `json:"displayName,omitempty"`
	Type        string    `json:"type,omitempty"`
}

func (m *Message) Valid() bool {
	if m.Message == "" {
		return false
	}
	if m.Room == "" {
		return false
	}
	if m.DisplayName == "" {
		return false
	}
	if m.Type == "" {
		return false
	}
	return true
}
