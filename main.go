package main

import (
	"log"
	"net/http"

	"github.com/snipextt/lets-talk-server/pkg"
)

func main() {
	http.HandleFunc("/ws", pkg.HandleWS)
	log.Println("Listining on *:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
