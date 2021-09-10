package internal

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Room struct {
	Users       []string
	Connections []*websocket.Conn
}

func (r *Room) Join(c *Client, message Message) {
	if c.isPartOfRoom(message.Room) {
		c.Conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Type:      "error",
			Message:   "Room already joined",
		})
		return
	}
	c.JoinedRooms = append(c.JoinedRooms, message.Room)
	r.Users = append(r.Users, message.DisplayName)
	r.Connections = append(r.Connections, c.Conn)
	log.Println(r)
	for _, user := range r.Users {
		c.Conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Type:      "UserAdd",
			Message:   user,
		})
	}
	for _, conn := range r.Connections {
		conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Type:      "UserAdd",
			Message:   message.DisplayName,
		})
	}

}
