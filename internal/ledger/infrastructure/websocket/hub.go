package websocket

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

type LedgerActivityHub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Mutex      sync.Mutex
}

func NewLedgerActivityHub() *LedgerActivityHub {
	return &LedgerActivityHub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *LedgerActivityHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			h.Clients[client] = true
			h.Mutex.Unlock()

		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Conn.Close()
			}
			h.Mutex.Unlock()

		case message := <-h.Broadcast:
			h.Mutex.Lock()
			for client := range h.Clients {
				if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
					client.Conn.Close()
					delete(h.Clients, client)
				}
			}
			h.Mutex.Unlock()
		}
	}
}

func (h *LedgerActivityHub) BroadcastMessage(message []byte) {
	h.Broadcast <- message
}
