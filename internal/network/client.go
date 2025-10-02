package network

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MessageType represents client and server message types
type ClientMessageType string
type ServerMessageType string

const (
	// Client -> Server
	ClientRegister      ClientMessageType = "REGISTER"
	ClientHeartbeat     ClientMessageType = "HEARTBEAT"
	ClientJoinRandom    ClientMessageType = "JOIN_RANDOM_MATCH"
	ClientJoinCode      ClientMessageType = "JOIN_CODE_MATCH"
	ClientCancelMatch   ClientMessageType = "CANCEL_MATCHING"
	ClientCall          ClientMessageType = "CALL"
	ClientRaise         ClientMessageType = "RAISE"
	ClientFold          ClientMessageType = "FOLD"
	ClientCheck         ClientMessageType = "CHECK"
	ClientAllIn         ClientMessageType = "ALL_IN"
	ClientLeaveGame     ClientMessageType = "LEAVE_GAME"
	ClientChatMessage   ClientMessageType = "CHAT_MESSAGE"

	// Server -> Client
	ServerRegisterSuccess ServerMessageType = "REGISTER_SUCCESS"
	ServerRegisterFailure ServerMessageType = "REGISTER_FAILURE"
	ServerMatchingStarted ServerMessageType = "MATCHING_STARTED"
	ServerMatchingProgress ServerMessageType = "MATCHING_PROGRESS"
	ServerMatchingCompleted ServerMessageType = "MATCHING_COMPLETED"
	ServerMatchingCancelled ServerMessageType = "MATCHING_CANCELLED"
	ServerGameStarted      ServerMessageType = "GAME_STARTED"
	ServerGameStateUpdate  ServerMessageType = "GAME_STATE_UPDATE"
	ServerPlayerAction     ServerMessageType = "PLAYER_ACTION"
	ServerRoundCompleted   ServerMessageType = "ROUND_COMPLETED"
	ServerGameEnded        ServerMessageType = "GAME_ENDED"
	ServerError            ServerMessageType = "ERROR"
	ServerInvalidAction    ServerMessageType = "INVALID_ACTION"
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
	conn       *websocket.Conn
	serverURL  string
	uuid       string
	nickname   string
	inbound    chan ServerMessage
	outbound   chan ClientMessage
	done       chan struct{}
	connected  bool
	mu         sync.RWMutex
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
		return c.conn.Close()
	}

	return nil
}

// IsConnected returns connection status
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			return
		}

		var serverMsg ServerMessage
		if err := json.Unmarshal(message, &serverMsg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
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
				log.Printf("Failed to marshal message: %v", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Failed to write message: %v", err)
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
