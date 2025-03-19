package main

import (
	"flag"
	"go-websocket-chats/server"
	"log"
)

func main(){
	// load .env 

	// run the server using net/http


	// create a new chat server
	// start the server
	addr:=flag.String("addr",":8080","HTTP service address") // get the address from config/env file
	flag.Parse()
	ChatServer:=server.NewChatServer(*addr)
	if err:=ChatServer.Start(); err !=nil{
		log.Fatal("server error:",err)
	}
}