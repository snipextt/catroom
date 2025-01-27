package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"
	"github.com/snipextt/catroom/pkg"
	"github.com/snipextt/catroom/pkg/db"
)

func main() {
	port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("captured %v, exiting gracefully...", sig)
			_, err := db.Client.FlushAll(context.Background()).Result()
			if err != nil {
				log.Fatal("Unable to clear database")
			}
			os.Exit(0)
		}
	}()
	db.ConnectClient()
	http.HandleFunc("/ws", pkg.HandleWS)
	http.HandleFunc("/", http.FileServer(http.Dir("./static")).ServeHTTP)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
