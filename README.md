# PokerHole Client

Modular Go client for PokerHole poker game using Bubble Tea TUI framework.

## Architecture

```
client/
├── cmd/
│   └── poker-client/
│       └── main.go                 # Entry point
│
├── internal/
│   ├── identity/
│   │   └── identity.go             # UUID and nickname management
│   │
│   ├── network/
│   │   └── client.go               # WebSocket client and protocol
│   │
│   ├── state/
│   │   └── game_state.go           # Game state management
│   │
│   └── ui/
│       └── model.go                # Bubble Tea UI model
│
├── go.mod
└── go.sum
```

## Building

```bash
go build -o poker-client cmd/poker-client/main.go
```

## Running

```bash
# Default (localhost:8080)
./poker-client

# Custom server
POKERHOLE_SERVER=ws://example.com:8080/ws/game ./poker-client
```

## Protocol

The client uses JSON-based WebSocket protocol with the following message types:

### Client -> Server
- `REGISTER` - Initial connection with UUID and nickname
- `HEARTBEAT` - Keep connection alive
- `JOIN_RANDOM_MATCH` - Join random matching
- `CALL`, `RAISE`, `FOLD`, `CHECK`, `ALL_IN` - Game actions

### Server -> Client
- `REGISTER_SUCCESS` - Registration confirmed
- `GAME_STATE_UPDATE` - Game state synchronization
- `PLAYER_ACTION` - Player action notifications
- `ERROR` - Error messages

## Features

- [x] Modular package structure
- [x] WebSocket client with automatic reconnection
- [x] JSON protocol support
- [x] State management for game synchronization
- [x] UUID-based player identification
- [x] Random nickname generation
- [x] Bubble Tea TUI integration
- [ ] Game action handling
- [ ] Matching system integration
- [ ] Chat functionality

## Legacy Files

- `main.go.legacy` - Original monolithic implementation
- `main_test.go.legacy` - Original tests
