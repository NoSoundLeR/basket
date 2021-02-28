package hub

// Hub ...
type Hub struct {
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
}

// Message ...
type Message struct {
	ID    string
	Value []byte
	Close bool
}

// NewHub ...
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

// BroadcastThrow ...
func (h *Hub) BroadcastThrow(id string, value string) {
	b := make([]byte, 0, len(value)+1)
	b = append(b, '+')
	b = append(b, []byte(value)...)
	h.Broadcast <- Message{
		ID:    id,
		Value: b,
		Close: false,
	}
}

// BroadcastResult ...
func (h *Hub) BroadcastResult(id string, value string) {
	b := make([]byte, 0, len(value)+1)
	b = append(b, '=')
	b = append(b, []byte(value)...)
	h.Broadcast <- Message{
		ID:    id,
		Value: b,
		Close: true,
	}
}

// Run ...
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				if client.ID == message.ID {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.Clients, client)
					}
				}
			}
		}
	}
}
