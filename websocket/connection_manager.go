package websocket

type ConnectionManager struct {
	ActiveClients                         map[*ClientHandler]bool
	BroadcastMessageToAllConnectedClients chan []byte
	AddNewConnectedClients                chan *ClientHandler
	RemoveDisconnectedClient              chan *ClientHandler
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		BroadcastMessageToAllConnectedClients: make(chan []byte),
		AddNewConnectedClients:                make(chan *ClientHandler),
		RemoveDisconnectedClient:              make(chan *ClientHandler),
		ActiveClients:                         make(map[*ClientHandler]bool),
	}
}

func (cm *ConnectionManager) Run() {
	for {
		select {
		// When a new client connects
		case client := <-cm.AddNewConnectedClients:
			cm.ActiveClients[client] = true
			// When a client disconnects
		case client := <-cm.RemoveDisconnectedClient:
			if _, ok := cm.ActiveClients[client]; ok {
				delete(cm.ActiveClients, client)
				close(client.send)
			}
			// When a message needs to be broadcast
		case message := <-cm.BroadcastMessageToAllConnectedClients:
			for client := range cm.ActiveClients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(cm.ActiveClients, client)
				}
			}
		}
	}
}
