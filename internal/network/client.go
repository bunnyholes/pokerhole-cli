package network

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MessageType represents client and server message types
type ClientMessageType string
type ServerMessageType string

const (
	// Client -> Server
	ClientRegister    ClientMessageType = "REGISTER"
	ClientHeartbeat   ClientMessageType = "HEARTBEAT"
	ClientJoinRandom  ClientMessageType = "JOIN_RANDOM_MATCH"
	ClientJoinCode    ClientMessageType = "JOIN_CODE_MATCH"
	ClientCancelMatch ClientMessageType = "CANCEL_MATCHING"
	ClientCall        ClientMessageType = "CALL"
	ClientRaise       ClientMessageType = "RAISE"
	ClientFold        ClientMessageType = "FOLD"
	ClientCheck       ClientMessageType = "CHECK"
	ClientAllIn       ClientMessageType = "ALL_IN"
	ClientLeaveGame   ClientMessageType = "LEAVE_GAME"
	ClientChatMessage ClientMessageType = "CHAT_MESSAGE"

	// Server -> Client
	ServerRegisterSuccess   ServerMessageType = "REGISTER_SUCCESS"
	ServerRegisterFailure   ServerMessageType = "REGISTER_FAILURE"
	ServerMatchingStarted   ServerMessageType = "MATCHING_STARTED"
	ServerMatchingProgress  ServerMessageType = "MATCHING_PROGRESS"
	ServerMatchingCompleted ServerMessageType = "MATCHING_COMPLETED"
	ServerMatchingCancelled ServerMessageType = "MATCHING_CANCELLED"
	ServerGameStarted       ServerMessageType = "GAME_STARTED"
	ServerGameStateUpdate   ServerMessageType = "GAME_STATE_UPDATE"
	ServerPlayerAction      ServerMessageType = "PLAYER_ACTION"
	ServerTurnChanged       ServerMessageType = "TURN_CHANGED"     // ADDED: Turn changed notification
	ServerRoundProgressed   ServerMessageType = "ROUND_PROGRESSED" // ADDED: Round progression (FLOP/TURN/RIVER)
	ServerRoundCompleted    ServerMessageType = "ROUND_COMPLETED"
	ServerGameEnded         ServerMessageType = "GAME_ENDED"
	ServerChatMessage       ServerMessageType = "CHAT_MESSAGE" // ADDED: Chat message
	ServerError             ServerMessageType = "ERROR"
	ServerInvalidAction     ServerMessageType = "INVALID_ACTION"
)

// ClientMessage represents a message from client to server
type ClientMessage struct {
	Type      ClientMessageType      `json:"type"`
	Timestamp int64                  `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

// ServerMessage represents a message from server to client
type ServerMessage struct {
	Type      ServerMessageType      `json:"type"`
	Timestamp int64                  `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

// Client represents a WebSocket client
type Client struct {
	conn      *websocket.Conn
	serverURL string
	uuid      string
	nickname  string
	inbound   chan ServerMessage
	outbound  chan ClientMessage
	done      chan struct{}
	connected bool
	mu        sync.RWMutex
}

// NewClient creates a new WebSocket client
func NewClient(serverURL, uuid, nickname string) *Client {
	return &Client{
		serverURL: serverURL,
		uuid:      uuid,
		nickname:  nickname,
		inbound:   make(chan ServerMessage, 100),
		outbound:  make(chan ClientMessage, 100),
		done:      make(chan struct{}),
	}
}

// Connect establishes WebSocket connection
func (c *Client) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.serverURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.connected = true
	c.mu.Unlock()

	// Start read/write goroutines
	go c.readPump()
	go c.writePump()

	// Send REGISTER message
	return c.Send(ClientMessage{
		Type:      ClientRegister,
		Timestamp: time.Now().UnixMilli(),
		Payload: map[string]interface{}{
			"uuid":     c.uuid,
			"nickname": c.nickname,
		},
	})
}

// ConnectWithTimeout establishes WebSocket connection with timeout
func (c *Client) ConnectWithTimeout(timeout time.Duration) error {
	// Create dialer with handshake timeout
	dialer := websocket.Dialer{
		HandshakeTimeout: timeout,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Dial with context
	conn, _, err := dialer.DialContext(ctx, c.serverURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect within %v: %w", timeout, err)
	}

	c.mu.Lock()
	c.conn = conn
	c.connected = true
	c.mu.Unlock()

	// Start read/write goroutines
	go c.readPump()
	go c.writePump()

	// Send REGISTER message
	return c.Send(ClientMessage{
		Type:      ClientRegister,
		Timestamp: time.Now().UnixMilli(),
		Payload: map[string]interface{}{
			"uuid":     c.uuid,
			"nickname": c.nickname,
		},
	})
}

// Send sends a message to the server
func (c *Client) Send(msg ClientMessage) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return fmt.Errorf("not connected")
	}

	select {
	case c.outbound <- msg:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("send timeout")
	}
}

// Receive returns the inbound message channel
func (c *Client) Receive() <-chan ServerMessage {
	return c.inbound
}

// Close closes the WebSocket connection
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil
	}

	c.connected = false
	close(c.done)

	if c.conn != nil {
		err := c.conn.Close()
		return err
	}

	return nil
}

// IsConnected returns connection status
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// JoinRandomMatch sends a request to join random match
func (c *Client) JoinRandomMatch() error {
	return c.Send(ClientMessage{
		Type:      ClientJoinRandom,
		Timestamp: time.Now().UnixMilli(),
		Payload:   map[string]interface{}{},
	})
}

// readPump reads messages from the WebSocket
func (c *Client) readPump() {
	defer func() {
		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		var serverMsg ServerMessage
		if err := json.Unmarshal(message, &serverMsg); err != nil {
			continue
		}

		select {
		case c.inbound <- serverMsg:
		case <-c.done:
			return
		}
	}
}

// writePump writes messages to the WebSocket
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-c.outbound:
			data, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-ticker.C:
			// Send heartbeat
			heartbeat := ClientMessage{
				Type:      ClientHeartbeat,
				Timestamp: time.Now().UnixMilli(),
			}
			data, _ := json.Marshal(heartbeat)
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-c.done:
			return
		}
	}
}

// SendGameAction sends a game action (FOLD, CHECK, CALL, RAISE, ALL_IN)
func (c *Client) SendGameAction(action ClientMessageType, amount int) error {
	payload := map[string]interface{}{}
	if action == ClientRaise && amount > 0 {
		payload["amount"] = amount
	}

	return c.Send(ClientMessage{
		Type:      action,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	})
}

// Fold sends FOLD action
func (c *Client) Fold() error {
	return c.SendGameAction(ClientFold, 0)
}

// Check sends CHECK action
func (c *Client) Check() error {
	return c.SendGameAction(ClientCheck, 0)
}

// Call sends CALL action
func (c *Client) Call() error {
	return c.SendGameAction(ClientCall, 0)
}

// Raise sends RAISE action with amount
func (c *Client) Raise(amount int) error {
	return c.SendGameAction(ClientRaise, amount)
}

// AllIn sends ALL_IN action
func (c *Client) AllIn() error {
	return c.SendGameAction(ClientAllIn, 0)
}
