package main

import (
	"log"
	"net/http"

	"github.com/snipextt/lets-talk-server/pkg"
)

func main() {
	// abc := time.Now()
	http.HandleFunc("/ws", pkg.HandleWS)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
