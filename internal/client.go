package internal

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var Rooms = make(map[string]Room)

type Client struct {
	Conn *websocket.Conn
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
			c.JoinRoom(message.Room, message.DisplayName)
		case "MESSAGE":
			c.SendMessage(message)
		case "CREATE":
			c.CreateRoom(message)
		}

	}
}

func (c *Client) JoinRoom(id string, displayName string) {
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
		for _, conn := range room.Connections {
			room.SendUserInRoomInfo(conn, displayName)
		}
		room.Connections = append(room.Connections, c.Conn)
		room.Users = append(room.Users, displayName)
		c.Conn.WriteJSON(&Message{
			Timestamp: time.Now(),
			Message:   "Joined room",
			Type:      "info",
		})
		room.SendUserInRoomInfo(c.Conn, nil)
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
	room := Room{
		Connections: []*websocket.Conn{c.Conn},
		Users:       []string{message.DisplayName},
	}
	Rooms[message.Room] = room
	c.JoinRoom(message.Room, message.DisplayName)
}
