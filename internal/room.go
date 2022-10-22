package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/snipextt/catroom/pkg/db"
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
			Message:   "You are already part of the room.",
		})
		return
	}
	for _, user := range r.Users {
		if user == message.DisplayName {
			c.Conn.WriteJSON(&Message{
				Timestamp: time.Now(),
				Type:      "error",
				Message:   fmt.Sprint("User with name: ", message.DisplayName, " already exists. Please use a different name to join the room: ", r.ID),
				Room:      r.ID,
			})
			return
		}
	}
	c.JoinedRooms = append(c.JoinedRooms, message.Room)
	c.DisplayNames = append(c.DisplayNames, message.DisplayName)
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
			Room:      message.Room,
		})
	}

}

func (r *Room) Leave(c *Client, name string) {
	for i, user := range r.Users {
		if user == name {
			r.Users = append(r.Users[:i], r.Users[i+1:]...)
			break
		}
	}
	for i, conn := range r.Connections {
		if conn == c.Conn {
			r.Connections = append(r.Connections[:i], r.Connections[i+1:]...)
			continue
		}
		conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Message:   name,
			Room:      r.ID,
			Type:      "UserRemove",
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
