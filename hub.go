package hub

import (
	"log"
	"net/http"
)

// TODO sync connections dictionary
// TODO rooms, i.e. namespaces
// TODO non JSON encoder

// Hub maintains the set of active connections and broadcasts messages to the connections.
type Hub struct {
	// Registered connections.
	connections map[*Conn]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *Conn

	// Unregister requests from connections.
	unregister chan *Conn
}

// NewHub creates new instance of Hub
func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Conn),
		unregister:  make(chan *Conn),
		connections: make(map[*Conn]bool),
	}
}

// Run starts the hub.
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &Conn{send: make(chan []byte, 256), ws: ws, hub: h}
	h.register <- c
	go c.writePump()
	c.readPump()
}

// Send broadcast message to all connections
func (h *Hub) Send(payload interface{}) {
	var msg = toJSON(payload)
	for c := range h.connections {
		c.send <- msg
	}
}
