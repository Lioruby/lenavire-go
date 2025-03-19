package adapters

import (
	"encoding/json"
	"lenavire/internal/ledger/infrastructure/websocket"
)

type WebSocketLedgerActivityChannel struct {
	hub *websocket.LedgerActivityHub
}

type LedgerActivityMessage struct {
	ShouldRefetch bool   `json:"shouldRefetch"`
	Type          string `json:"type"`
}

func NewWebSocketLedgerActivityChannel(hub *websocket.LedgerActivityHub) *WebSocketLedgerActivityChannel {
	return &WebSocketLedgerActivityChannel{
		hub: hub,
	}
}

func (c *WebSocketLedgerActivityChannel) Send(messageType string) error {
	message := LedgerActivityMessage{
		ShouldRefetch: true,
		Type:          messageType,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	c.hub.BroadcastMessage(jsonMessage)
	return nil
}
