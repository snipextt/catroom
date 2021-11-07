package internal

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/snipextt/lets-talk-server/pkg/db"
)

type Room struct {
	Users       []string
	Connections []*websocket.Conn
	Messages    chan Message
	ID          string
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
	for _, user := range r.Users {
		c.Conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Type:      "UserAdd",
			Message:   user,
			Room:      message.Room,
		})
	}
	val, err := db.Client.LRange(context.Background(), r.ID, 0, -1).Result()
	if err != nil {
		log.Println(err.Error())
	}
	c.Conn.WriteJSON(&MessageHistory{
		Timestamp: time.Now(),
		Room:      r.ID,
		Type:      "History",
		History:   val,
	})
	for _, conn := range r.Connections {
		if conn == c.Conn {
			continue
		}
		conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Type:      "UserAdd",
			Message:   message.DisplayName,
		})
	}

}

func (r *Room) WatchMessages() {
	for {
		msg := <-r.Messages
		p, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
		}
		db.Client.LPush(context.Background(), r.ID, p)
	}
}
