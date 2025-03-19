package main

import (
	"flag"
	"go-websocket-chats/server"
	"log"
)

func main(){
	addr:=flag.String("addr",":8080","HTTP service address")
	flag.Parse()
	ChatServer:=server.NewChatServer(*addr)
	if err:=ChatServer.Start(); err !=nil{
		log.Fatal("server error:",err)
	}
}