package pkg

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/snipextt/catroom/internal"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWS(rw http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println("Cannot switch protocals!")
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ws.WriteMessage(websocket.TextMessage, []byte("Connected!"))
	clientInstance := &internal.Client{
		Conn: ws,
		// brodcastedMessages: make(chan Message, 100),
	}
	go clientInstance.WatchMessages()
}
