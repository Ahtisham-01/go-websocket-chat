package websocket

type ConnectionManager struct {
	ActiveClients                         map[*ClientHandler]bool
	BroadcastMessageToAllConnectedClients chan []byte
	AddNewConnectedClients                chan *ClientHandler
	RemoveDisconnectedClient                chan *ClientHandler
}

func NewConnectionManager() *ConnectionManager{
	return &ConnectionManager{
		BroadcastMessageToAllConnectedClients: make(chan []byte),
		AddNewConnectedClients: make(chan *ClientHandler),
		RemoveDisconnectedClient: make(chan *ClientHandler),
		ActiveClients: make(map[*ClientHandler]bool),
	}
}

