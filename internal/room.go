package internal

import "github.com/gorilla/websocket"

type Room struct {
	Users       []string
	Connections []*websocket.Conn
}

func (r *Room) SendUserInRoomInfo(c *websocket.Conn, user interface{}) {
	if user != "" && user != nil {
		c.WriteJSON(&Message{
			Type:    "UserAdd",
			Message: user.(string),
		})
	} else {
		for _, user := range r.Users {
			c.WriteJSON(&Message{
				Type:    "UserAdd",
				Message: user,
			})
		}
	}
}
