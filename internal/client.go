package internal

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var Rooms = make(map[string][]*websocket.Conn)

type Client struct {
	Conn *websocket.Conn
	// brodcastedMessages chan Message
}

func (c *Client) WatchMessages() {
	for {
		var message Message
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			log.Println(err.Error())
			c.Conn.Close()
			return
		}
		switch message.Type {
		case "JOIN":
			c.JoinRoom(message.Room)
		case "MESSAGE":
			c.SendMessage(message)
		}

	}
}

func (c *Client) JoinRoom(id string) {
	if id == "" {
		c.Conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Message:   "Room id is missing",
			Type:      "error",
		})
		return
	}
	if _, ok := Rooms[id]; !ok {
		c.Conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Message:   "Room id invalid",
			Type:      "error",
		})
		return
	} else {
		Rooms[id] = append(Rooms[id], c.Conn)
		c.Conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Message:   "Joined room",
			Type:      "info",
		})
	}
}

func (c *Client) SendMessage(message Message) {

}
