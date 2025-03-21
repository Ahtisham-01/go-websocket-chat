package websocket

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	totalTimeToWriteMessage = 10 * time.Second
	TotalTimeToGetResponse  = 60 * time.Second
	pingPeriod              = (TotalTimeToGetResponse * 9) / 10
	maxMessageSize          = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var webSocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type ClientHandler struct {
	manager *ConnectionManager
	conn    *websocket.Conn
	send    chan []byte
}

func (c *ClientHandler) readMessages() {

	defer func() {
		c.manager.RemoveDisconnectedClient <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(TotalTimeToGetResponse))

	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(TotalTimeToGetResponse))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		c.manager.BroadcastMessageToAllConnectedClients <- message

	}
}

func (c *ClientHandler) writeMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(TotalTimeToGetResponse))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(totalTimeToWriteMessage))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func HandleWebSocketConnection(manager *ConnectionManager, w http.ResponseWriter, r *http.Request) {
	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket Upgrade error:", err)
		return
	}
	client := &ClientHandler{
		manager: manager,
		conn:    conn,
		send:    make(chan []byte, 256),
	}
	client.manager.AddNewConnectedClients <- client
	go client.writeMessages()
	go client.readMessages()
}
