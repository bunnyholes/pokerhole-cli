# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**pokerhole-cli** is a Go-based terminal client for the PokerHole Texas Hold'em poker game. It communicates with the backend server via WebSocket and provides a modern TUI using the Bubble Tea framework.

**Design Goal**: Support both **network mode** (WebSocket to backend server) and **standalone mode** (local game with AI players). Currently only network mode is implemented.

## Build and Run Commands

```bash
# Build the client
go build -o poker-client cmd/poker-client/main.go

# Run with default server (localhost:8080)
./poker-client

# Run with custom server
POKERHOLE_SERVER=ws://example.com:8080/ws/game ./poker-client

# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test -v ./internal/network

# Download dependencies
go mod download

# Update dependencies
go mod tidy
```

## Cross-Platform Builds

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o poker-client-linux cmd/poker-client/main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o poker-client-macos-intel cmd/poker-client/main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o poker-client-macos-arm cmd/poker-client/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o poker-client.exe cmd/poker-client/main.go
```

## Architecture Overview

This is a **Go-based terminal UI client** for the PokerHole poker game. It connects to a Java/Spring Boot backend server via WebSocket.

### Current Package Structure (Network Mode Only)

```
cmd/poker-client/main.go          # Entry point, initializes and wires components
internal/identity/identity.go     # UUID persistence (~/.pokerhole/uuid) and nickname generation
internal/network/client.go        # WebSocket client with JSON protocol, heartbeat (30s)
internal/state/game_state.go      # Thread-safe game state (RWMutex protected)
internal/ui/model.go              # Bubble Tea TUI model (View/Update/Init pattern)
```

### Data Flow

1. **Initialization**: `main.go` → creates WebSocket client → connects to server
2. **Registration**: Client sends `REGISTER` message with UUID + nickname
3. **Message Loop**:
   - `network.Client.readPump()` receives messages → pushes to `inbound` channel
   - `ui.listenForMessages()` subscribes to channel → sends to Bubble Tea
   - `Model.Update()` processes messages → updates view
4. **Heartbeat**: `network.Client.writePump()` sends heartbeat every 30s
5. **State Management**: `state.GameState` stores game state with mutex protection

### Key Design Patterns

**WebSocket Client (`internal/network/client.go`)**:
- Goroutine-based read/write pumps (concurrent message handling)
- Channel-based inbound/outbound queues (buffered 100 messages)
- Thread-safe connection status tracking with `sync.RWMutex`
- No automatic reconnection (manual restart required on disconnect)

**Bubble Tea UI (`internal/ui/model.go`)**:
- Model-View-Update pattern (functional reactive programming)
- View modes: `ViewSplash` → `ViewConnecting` → `ViewMenu` → `ViewGame`
- Custom message types: `ServerMessageMsg`, `ConnectionEstablishedMsg`
- Commands return `tea.Cmd` for async operations

**State Synchronization (`internal/state/game_state.go`)**:
- `Update()`: Parses server payload and updates internal state (write lock)
- `GetSnapshot()`: Returns immutable copy for rendering (read lock)
- Thread-safe for concurrent UI reads and network updates

**UUID Persistence**:
- Stored in `~/.pokerhole/uuid` file
- Created on first run if missing
- Ensures consistent player identity across sessions

### Component Dependencies

```
main.go
  ├─> identity.GetOrCreateUUID()      # Load/create persistent UUID
  ├─> identity.GenerateNickname()     # Random name like "LuckyShark123"
  ├─> network.NewClient()             # WebSocket client
  └─> ui.NewModel()                   # Bubble Tea UI
      └─> uses network.Client for messaging
```

### Server Integration

**Expected Server Endpoint:** `ws://localhost:8080/ws/game` (configurable via `POKERHOLE_SERVER`)

**Initial Handshake:**
```
Client sends: {"type":"REGISTER","timestamp":...,"payload":{"uuid":"...","nickname":"..."}}
Server responds: {"type":"REGISTER_SUCCESS","timestamp":...,"payload":{...}}
```

**Related Repository:** [PokerHole Server](https://github.com/bunnyholes/pokerhole) (Java/Spring Boot)

## Development Notes

### Adding New Message Types

1. Add constant to `internal/network/client.go` (e.g., `ClientNewAction ClientMessageType = "NEW_ACTION"`)
2. Update protocol documentation in README.md
3. Handle in `internal/ui/model.go` Update() method for server messages

### Bubble Tea Integration

- `internal/ui/model.go` implements `tea.Model` interface
- `Update()` handles messages (keyboard, server messages)
- `View()` renders TUI
- Use `tea.WithAltScreen()` for full-screen TUI

### Logging

- Logs written to `/tmp/pokerhole-client.log`
- Use `log.Printf()` for debugging
- Check log file if client behaves unexpectedly

### Thread Safety

- **ALWAYS** use mutex when accessing `Client.connected` or `Client.conn`
- State reads/writes in `game_state.go` must use RLock/Lock
- Channels are already thread-safe (no mutex needed)

## Testing the Client

1. Start PokerHole server: `cd ../pokerhole && ./gradlew bootRun`
2. Build client: `go build -o poker-client cmd/poker-client/main.go`
3. Run client: `./poker-client`
4. Check logs: `tail -f /tmp/pokerhole-client.log`
5. Verify WebSocket connection in server logs
