package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var Rooms = make(map[string]*Room)

type Client struct {
	Conn        *websocket.Conn
	JoinedRooms []string
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
		if !message.Valid() {
			c.Conn.WriteJSON(&Message{
				Timestamp: time.Now(),
				Message:   "Invalid message syntax!",
				Type:      "error",
			})
			continue
		}
		switch message.Type {
		case "JOIN":
			c.JoinRoom(message)
		case "MESSAGE":
			c.SendMessage(message)
		case "CREATE":
			c.CreateRoom(message)
		}
		fmt.Println(Rooms)
	}
}

func (c *Client) JoinRoom(message Message) {
	id := message.Room
	if id == "" {
		c.Conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Message:   "Room id is missing",
			Type:      "error",
		})
		return
	}
	if room, ok := Rooms[id]; !ok {
		c.Conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Message:   "Room id invalid",
			Type:      "error",
		})
		return
	} else {
		room.Join(c, message)
	}
}

func (c *Client) SendMessage(message Message) {
	if message.Room == "" {
		c.Conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Message:   "Room id is missing",
			Type:      "error",
		})
		return
	}
	if room, ok := Rooms[message.Room]; !ok {
		c.Conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Message:   "Room id invalid",
			Type:      "error",
		})
		return
	} else {
		for _, conn := range room.Connections {
			conn.WriteJSON(&message)
		}
	}
}

func (c *Client) CreateRoom(message Message) {
	if message.Room == "" {
		c.Conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Message:   "Room id is missing",
			Type:      "error",
		})
		return
	}
	if _, ok := Rooms[message.Room]; ok {
		c.Conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Message:   "Room id already exists",
			Type:      "error",
		})
		return
	}
	room := &Room{}
	Rooms[message.Room] = room
	c.JoinRoom(message)
}

func (c *Client) isPartOfRoom(roomToJoin string) bool {
	for _, room := range c.JoinedRooms {
		if room == roomToJoin {
			return true
		}
	}
	return false
}
