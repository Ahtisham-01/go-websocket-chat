package main

import (
	"log"
	"os"
	"go-websocket-chats/server"
)

// init function to load environment variables
func init() {
	file, err := os.Open(".env")
	if err != nil {
		log.Println("Warning: No .env file found")
		return
	}
	defer file.Close()
}

func main() {
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080" // fallback
	}
	ChatServer := server.NewChatServer(addr)
	// Start the server
	if err := ChatServer.Start(); err != nil {
		log.Fatal("server error:", err)
	}
}
