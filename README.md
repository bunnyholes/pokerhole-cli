# PokerHole CLI

Go terminal client for Texas Hold'em poker with offline support and event sourcing.

---

## Quick Start

```bash
# Run with default server (localhost:8080)
go run cmd/poker-client/main.go

# Run with custom server
POKERHOLE_SERVER=ws://example.com:8080/ws/game go run cmd/poker-client/main.go

# Build executable
go build -o poker cmd/poker-client/main.go
./poker
```

---

## Architecture

### Hexagonal Architecture + Event Sourcing + Offline-First

```
┌───────────────────────────────────────────────┐
│              TUI Layer                        │
│        Bubble Tea + Lipgloss                  │
│    (Splash → Menu → Matching → Game)         │
└──────────────────┬────────────────────────────┘
                   │
                   ↓
┌───────────────────────────────────────────────┐
│          Application Layer                    │
│      Use Cases | Commands | Queries          │
│           Port Interfaces                     │
└──────────────────┬────────────────────────────┘
                   │
                   ↓
┌───────────────────────────────────────────────┐
│            Domain Layer                       │
│         Pure Go (no framework deps)           │
│    Game | Player | Card | HandEvaluator      │
│    (Ported from Java, golden-tested)         │
└──────────────────┬────────────────────────────┘
                   │
                   ↓
┌───────────────────────────────────────────────┐
│           Output Adapters                     │
│  SQLite Event Store | WebSocket | Crypto     │
│      (ed25519 signatures for events)         │
└───────────────────────────────────────────────┘
```

### Architecture Decisions

See [`../docs/adr/`](../docs/adr/) for rationale:
- [ADR-001: Event Sourcing](../docs/adr/001-event-sourcing-for-gameplay.md) - Local SQLite event store
- [ADR-003: ed25519 Signatures](../docs/adr/003-ed25519-client-signatures.md) - Sign events for offline→online sync
- [ADR-005: Deterministic RNG](../docs/adr/005-deterministic-rng-fairness.md) - Same shuffle as Java server

---

## Project Structure

```
pokerhole-cli/
├── cmd/
│   └── poker-client/
│       └── main.go          # Entry point, wire components
│
├── internal/
│   ├── domain/              # Pure Go domain (ported from Java)
│   │   ├── game/           # Game aggregate, rules
│   │   │   ├── game.go
│   │   │   ├── round.go
│   │   │   └── event.go    # GameStarted, PlayerActed, etc.
│   │   ├── player/         # Player aggregate
│   │   ├── card/           # Card, Deck (deterministic shuffle)
│   │   └── evaluator/      # HandEvaluator (golden-tested)
│   │
│   ├── application/        # Use cases, ports
│   │   ├── port/
│   │   │   ├── in/         # StartGameUseCase, PlaceBetUseCase
│   │   │   └── out/        # GameRepositoryPort, EventStorePort
│   │   └── service/        # Use case implementations
│   │
│   ├── adapter/
│   │   ├── ui/             # Bubble Tea TUI
│   │   │   ├── model.go    # Main UI model
│   │   │   ├── splash.go   # Splash screen
│   │   │   ├── menu.go     # Menu screen
│   │   │   ├── game.go     # Game screen
│   │   │   └── styles.go   # Lipgloss styles
│   │   ├── storage/        # SQLite event store
│   │   │   └── sqlite.go
│   │   ├── network/        # WebSocket client
│   │   │   └── client.go   # Auto-reconnect, heartbeat
│   │   └── crypto/         # ed25519 signing
│   │       └── signer.go
│   │
│   └── identity/           # UUID persistence, nickname gen
│       └── identity.go
│
├── tests/
│   └── golden/             # Golden test vectors
│       ├── hand_eval_test.go
│       └── data/
│           └── hand_eval.json  # 21 test cases (Java parity)
│
├── go.mod
└── go.sum
```

---

## Features

### Current (Phase 2 Complete)

- ✅ **Modern TUI**: Bubble Tea framework with Lipgloss styling
- ✅ **WebSocket Client**: Real-time communication with server
- ✅ **Event Sourcing**: Local SQLite event store
- ✅ **Offline Support**: Play without server (AI opponents)
- ✅ **Go-Java Parity**: 21 golden tests ensure same behavior
- ✅ **ed25519 Signatures**: Sign events for security
- ✅ **UUID Persistence**: Consistent player identity
- ✅ **Auto-Reconnect**: Handles disconnections gracefully

### Planned (Phase 3+)

- ⚠️ **Online Multiplayer**: Play with others via server
- ⚠️ **Offline AI**: Conservative/Aggressive strategies
- ⚠️ **Offline→Online Sync**: Sync local games to server
- ⚠️ **Tournament Mode**: Multi-table tournaments
- ⚠️ **Replay Mode**: Replay past games from event log

---

## Tech Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Language | Go | 1.25+ |
| TUI Framework | Bubble Tea | latest |
| Styling | Lipgloss | latest |
| Database | SQLite | 3.x |
| Encryption | SQLCipher | optional |
| Crypto | ed25519 | stdlib |
| WebSocket | gorilla/websocket | latest |
| Testing | Go testing | stdlib |

---

## Domain Model (Go Port)

### Ported from Java

The Go domain model is a **direct port** of the Java server domain:

```go
// game/game.go
type Game struct {
    ID          GameID
    Players     []*Player
    Deck        *Deck
    Community   []Card
    Round       BettingRound
    Pot         int
    CurrentBet  int
    // ...
}

func (g *Game) StartGame(seed int64) []DomainEvent { /* ... */ }
func (g *Game) PlayerAction(playerID PlayerID, action PlayerAction, amount int) []DomainEvent { /* ... */ }
func (g *Game) ProgressRound() []DomainEvent { /* ... */ }
```

### Golden Test Verification

**21 test cases** ensure Go ↔ Java parity:

```bash
go test -v ./tests/golden/
```

**Test File**: `tests/golden/data/hand_eval.json` (same as server)

**Verified**:
- Hand evaluation (Royal Flush → High Card)
- Same inputs → Same outputs (deterministic)
- No divergence between Java and Go

---

## Event Sourcing

### Local Event Store (SQLite)

All game actions stored as events:

```sql
CREATE TABLE events (
    event_id TEXT PRIMARY KEY,
    game_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    payload TEXT NOT NULL,      -- JSON
    timestamp INTEGER NOT NULL,
    signature BLOB,              -- ed25519 signature
    synced BOOLEAN DEFAULT 0
);
```

### Event Types

```go
type DomainEvent interface {
    EventID() string
    GameID() string
    OccurredAt() time.Time
}

type GameStarted struct {
    EventID    string
    GameID     string
    OccurredAt time.Time
    PlayerIDs  []PlayerID
    Seed       int64
}

type PlayerActed struct {
    EventID    string
    GameID     string
    OccurredAt time.Time
    PlayerID   PlayerID
    Action     PlayerAction
    Amount     int
}
// ... etc
```

### Event Signatures (ed25519)

**All events signed** for security and offline→online sync:

```go
// Generate keypair
publicKey, privateKey, err := ed25519.GenerateKey(nil)

// Sign event
eventJSON := json.Marshal(event)
signature := ed25519.Sign(privateKey, eventJSON)

// Verify (server-side)
valid := ed25519.Verify(publicKey, eventJSON, signature)
```

**Benefits**:
- Tamper detection
- Non-repudiation (player can't deny action)
- Offline→online sync (server trusts signed events)

---

## TUI (Bubble Tea)

### View Modes

```
ViewSplash → ViewMenu → ViewMatching → ViewGame → ViewResults
```

**Splash Screen** (3 seconds):
```
  ╔═══════════════════════════════╗
  ║                               ║
  ║      🂡  POKERHOLE  🂱         ║
  ║                               ║
  ║    Texas Hold'em Poker        ║
  ║                               ║
  ╚═══════════════════════════════╝
```

**Menu Screen**:
```
  1. 🎲 Random Match (online)
  2. 🔢 Code Match (online)
  3. 🤖 Play vs AI (offline)
  4. 📜 How to Play
  5. ℹ️  About

  Q: Quit
```

**Game Screen**:
```
  Community Cards: [AS] [KH] [QD]

  Pot: $1,000 | Current Bet: $200

  Players:
  • LuckyShark123 (you)    $8,500  [JC] [TD]
  • Player2 (turn)         $9,200  [??] [??]
  • Player3                $7,800  [??] [??]
  • Player4                $10,500 [??] [??]

  Your turn!
  [C]all $200 | [R]aise | [F]old | [A]ll-in
```

### Keyboard Controls

| Key | Action |
|-----|--------|
| Arrow Keys / Tab | Navigate |
| Enter | Select |
| C / c | Call |
| R / r | Raise (prompts for amount) |
| F / f | Fold |
| A / a | All-in |
| H / h | Help |
| Q / Esc | Quit / Back |

---

## WebSocket Protocol

### Connection

```
ws://localhost:8080/ws/game
```

### Client → Server Messages

```go
type ClientMessage struct {
    Type      string      `json:"type"`
    Timestamp int64       `json:"timestamp"`
    Payload   interface{} `json:"payload"`
}

// Types
const (
    ClientRegister        = "REGISTER"         // {uuid, nickname}
    ClientHeartbeat       = "HEARTBEAT"        // {}
    ClientJoinRandom      = "JOIN_RANDOM_MATCH" // {}
    ClientCall            = "CALL"             // {}
    ClientRaise           = "RAISE"            // {amount}
    ClientFold            = "FOLD"             // {}
    // ...
)
```

### Server → Client Messages

```go
type ServerMessage struct {
    Type      string          `json:"type"`
    Timestamp int64           `json:"timestamp"`
    Payload   json.RawMessage `json:"payload"`
}

// Types
const (
    ServerRegisterSuccess = "REGISTER_SUCCESS" // {playerId}
    ServerGameStateUpdate = "GAME_STATE_UPDATE" // {full state}
    ServerPlayerAction    = "PLAYER_ACTION"     // {playerId, action}
    ServerError           = "ERROR"             // {message}
    // ...
)
```

### Auto-Reconnect

```go
func (c *Client) Connect() error {
    for {
        err := c.dial()
        if err == nil {
            return nil
        }
        log.Printf("Connection failed, retrying in 5s: %v", err)
        time.Sleep(5 * time.Second)
    }
}
```

**Heartbeat** (30s interval):
```go
func (c *Client) writePump() {
    ticker := time.NewTicker(30 * time.Second)
    for {
        select {
        case <-ticker.C:
            c.SendMessage(ClientHeartbeat, nil)
        }
    }
}
```

---

## Configuration

### Environment Variables

```bash
# Server URL
export POKERHOLE_SERVER=ws://localhost:8080/ws/game

# Data directory (default: ~/.pokerhole)
export POKERHOLE_DATA_DIR=/path/to/data

# Enable SQLCipher encryption
export POKERHOLE_ENCRYPT_DB=true

# Log level
export POKERHOLE_LOG_LEVEL=debug
```

### User Data

**Default location**: `~/.pokerhole/`

```
~/.pokerhole/
├── uuid               # Player UUID (persistent)
├── poker.db          # SQLite event store
├── keychain.db       # ed25519 keys (encrypted)
└── config.toml       # User preferences (future)
```

---

## Development

### Prerequisites

- **Go 1.25+**
- **SQLite3** (included in Go stdlib)

### Setup

```bash
# Clone repo
git clone <repo-url>
cd pokerhole-cli

# Download dependencies
go mod download

# Run
go run cmd/poker-client/main.go
```

### Hot Reload

Use `air` for hot reload during development:

```bash
# Install air
go install github.com/air-verse/air@latest

# Run with hot reload
air
```

`air.toml`:
```toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/poker cmd/poker-client/main.go"
bin = "tmp/poker"
include_ext = ["go"]
exclude_dir = ["tmp"]
```

### Debugging

```bash
# Run with Delve debugger
dlv debug cmd/poker-client/main.go

# Attach to running process
dlv attach $(pgrep poker)
```

---

## Testing

### Unit Tests

```bash
# All tests
go test ./...

# Verbose
go test -v ./...

# Specific package
go test -v ./internal/domain/card
```

### Golden Tests (Go ↔ Java Parity)

```bash
# Run golden tests
go test -v ./tests/golden/

# Update golden data (if Java changes)
cd ../pokerhole-server
./gradlew test --tests "*GoldenVectorValidationTest*"
# Copy updated hand_eval.json to CLI
```

**21 test cases** verify:
- Hand evaluation correctness
- Deterministic shuffling
- Same behavior as Java server

### Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View in browser
go tool cover -html=coverage.out
```

**Target**: 80%+ coverage (domain: 95%+)

---

## Building

### Development Build

```bash
go build -o poker cmd/poker-client/main.go
./poker
```

### Production Build

```bash
# With optimizations
go build -ldflags="-s -w" -o poker cmd/poker-client/main.go

# Check size
ls -lh poker
```

### Cross-Platform

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o poker-linux cmd/poker-client/main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o poker-macos-intel cmd/poker-client/main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o poker-macos-arm cmd/poker-client/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o poker.exe cmd/poker-client/main.go
```

---

## Performance

### Targets

| Metric | Target | Actual |
|--------|--------|--------|
| Hand evaluation | < 10ms | ✅ ~1ms |
| Event append | < 5ms | ✅ ~2ms |
| TUI render | < 16ms (60fps) | ✅ ~10ms |
| SQLite query | < 1ms | ✅ ~0.5ms |
| Startup time | < 1s | ✅ ~300ms |

### Profiling

```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

---

## Contributing

### Coding Standards

#### Naming

| Type | Convention | Example |
|------|-----------|---------|
| Interface | Noun or `*er` | `GameRepository`, `HandEvaluator` |
| Struct | Noun | `Game`, `Player`, `Card` |
| Function | Verb | `StartGame()`, `PlayerAction()` |
| Test | `Test*` | `TestHandEvaluator_RoyalFlush` |

#### Imports

Group imports:
```go
import (
    // stdlib
    "context"
    "fmt"

    // external
    tea "github.com/charmbracelet/bubbletea"

    // internal
    "github.com/bunnyholes/pokerhole-cli/internal/domain/game"
)
```

#### Error Handling

```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to start game: %w", err)
}

// Check errors explicitly (no silent failures)
// Use errors.Is() and errors.As() for comparison
```

### Pull Request Checklist

- [ ] All tests pass (`go test ./...`)
- [ ] Golden tests pass (Go ↔ Java parity)
- [ ] `go fmt` applied
- [ ] `go vet` clean
- [ ] `golangci-lint run` clean (no errors)
- [ ] Coverage ≥ 80% (domain ≥ 95%)
- [ ] No breaking changes to event schema

---

## Troubleshooting

### Database Locked

**Cause**: Another process has SQLite open

**Fix**:
```bash
rm ~/.pokerhole/poker.db
# Restart client
```

### UUID Not Persisting

**Check permissions**:
```bash
ls -la ~/.pokerhole/
chmod 600 ~/.pokerhole/uuid
```

### WebSocket Connection Failed

**Verify server**:
```bash
curl http://localhost:8080/actuator/health

# Test WebSocket
wscat -c ws://localhost:8080/ws/game
```

### TUI Not Rendering Correctly

**Check terminal**:
```bash
echo $TERM
# Should be: xterm-256color or similar

# Fix
export TERM=xterm-256color
```

### Golden Tests Failing

**Sync with Java**:
```bash
cd ../pokerhole-server
./gradlew test --tests "*GoldenVectorValidationTest*"

# Copy JSON
cp pokerhole-server/src/test/resources/golden/hand_eval.json \
   pokerhole-cli/tests/golden/data/
```

---

## Related Projects

- **[PokerHole Server](../pokerhole-server/)** - Java/Spring Boot backend
- **[Root Project](../)** - Monorepo root

---

## Documentation

- **Architecture**: See [`../docs/adr/`](../docs/adr/) for all architecture decisions
- **Protocol**: See server README for WebSocket protocol details
- **Golden Tests**: See `tests/golden/README.md`

---

## License

MIT License

---

## Status

**Last Updated**: 2025-10-03

- **Phase**: 2 complete (domain model, event store, TUI)
- **Go-Java Parity**: ✅ 21 golden tests passing
- **Production**: Not ready (Phase 3+ needed)

**Next Milestone**: Complete online multiplayer (Phase 3)
