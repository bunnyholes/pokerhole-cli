# PokerHole CLI

Terminal-based client for PokerHole Texas Hold'em poker game, built with Go and Bubble Tea TUI framework.

## Features

- Modern TUI interface with Bubble Tea framework
- WebSocket-based real-time communication
- Automatic reconnection handling
- Persistent UUID-based player identification
- Random nickname generation
- Thread-safe state management
- Modular package architecture

## Prerequisites

- Go 1.21 or higher
- PokerHole server running (default: `ws://localhost:8080/ws/game`)

## Installation

```bash
# Clone the repository
git clone https://github.com/bunnyholes/pokerhole-cli.git
cd pokerhole-cli

# Download dependencies
go mod download

# Build
go build -o poker-client cmd/poker-client/main.go
```

## Usage

```bash
# Run with default server (localhost:8080)
./poker-client

# Run with custom server
POKERHOLE_SERVER=ws://example.com:8080/ws/game ./poker-client
```

## Architecture

```
pokerhole-cli/
├── cmd/
│   └── poker-client/
│       └── main.go              # Entry point
│
├── internal/
│   ├── identity/
│   │   └── identity.go          # UUID persistence, nickname generation
│   │
│   ├── network/
│   │   └── client.go            # WebSocket client with auto-reconnect
│   │
│   ├── state/
│   │   └── game_state.go        # Thread-safe game state management
│   │
│   └── ui/
│       └── model.go             # Bubble Tea UI components
│
├── go.mod
└── go.sum
```

## WebSocket Protocol

### Client to Server

```json
{
  "type": "REGISTER",
  "timestamp": 1234567890,
  "payload": {
    "uuid": "player-uuid",
    "nickname": "PlayerName"
  }
}
```

**Message Types:**
- `REGISTER` - Initial connection
- `HEARTBEAT` - Keep-alive (auto-sent every 30s)
- `JOIN_RANDOM_MATCH` - Join random matching
- `JOIN_CODE_MATCH` - Join with matching code
- `CALL`, `RAISE`, `FOLD`, `CHECK`, `ALL_IN` - Game actions

### Server to Client

```json
{
  "type": "GAME_STATE_UPDATE",
  "timestamp": 1234567890,
  "payload": {
    "gameId": "game-123",
    "round": "FLOP",
    "pot": 1000,
    "currentBet": 200,
    "communityCards": ["AS", "KH", "QD"],
    "players": [...],
    "currentPlayer": "player-2",
    "validActions": ["CALL", "RAISE", "FOLD"]
  }
}
```

**Message Types:**
- `REGISTER_SUCCESS` - Registration confirmed
- `GAME_STATE_UPDATE` - Game state sync
- `PLAYER_ACTION` - Action notifications
- `MATCHING_STARTED`, `MATCHING_COMPLETED` - Matching events
- `ERROR` - Error messages

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o poker-client-linux cmd/poker-client/main.go

# macOS
GOOS=darwin GOARCH=arm64 go build -o poker-client-macos cmd/poker-client/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o poker-client.exe cmd/poker-client/main.go
```

## Configuration

### Environment Variables

- `POKERHOLE_SERVER` - WebSocket server URL (default: `ws://localhost:8080/ws/game`)

### User Data

- UUID stored in: `~/.pokerhole/uuid`
- Data directory created automatically on first run

## Features Status

- [x] Modular package structure
- [x] WebSocket client with automatic reconnection
- [x] JSON protocol support
- [x] Heartbeat mechanism (30s interval)
- [x] Thread-safe state management
- [x] UUID persistence
- [x] Random nickname generation
- [x] Bubble Tea TUI framework
- [ ] Complete game action handling
- [ ] Matching system UI
- [ ] Game table rendering
- [ ] Chat functionality
- [ ] Spectator mode

## Related Projects

- [PokerHole Server](https://github.com/bunnyholes/pokerhole) - Java/Spring Boot backend

## Contributing

See the main [PokerHole CONTRIBUTING.md](https://github.com/bunnyholes/pokerhole/blob/main/CONTRIBUTING.md) for development guidelines.

## License

MIT License

## Authors

- [@xiyo](https://github.com/xiyo)

## Legacy Code

The `*.legacy` files contain the original monolithic implementation, preserved for reference.
