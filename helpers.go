package hub

import (
	"encoding/json"
	"log"
	"net/http"
)

func toJSON(payload interface{}) []byte {
	var b, err = json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return b
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &Conn{send: make(chan []byte, 256), ws: ws, hub: hub}
	hub.register <- c
	go c.writePump()
	c.readPump()
}
