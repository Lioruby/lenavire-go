# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Running
- **Development with hot reload**: `air -c .air.toml` (uses Air for hot reloading)
- **Manual build**: `go build -o ./tmp/main ./cmd/api/main.go`
- **Run built binary**: `./tmp/main`
- **Docker development**: `docker-compose up` (includes PostgreSQL database)

### Testing
- **Run tests**: `go test ./...` 
- **Run specific test**: `go test ./internal/ledger/application/commands/add_expense/`
- **Test framework**: Uses Ginkgo/Gomega BDD testing framework

### Database
- **Local PostgreSQL**: Available via docker-compose on port 5432
- **Credentials**: admin/azerty, database: lenavire_dev
- **Auto-migration**: Handled automatically on startup via GORM

## Architecture

### Clean Architecture Pattern
The codebase follows Domain-Driven Design with Clean Architecture:

```
internal/ledger/
├── domain/           # Core business logic
│   ├── entities/     # Business entities (Expense, Payment)
│   ├── valuesobjects/# Value objects (Amount, PaymentType)
│   └── exceptions/   # Domain exceptions
├── application/      # Use cases and business rules
│   ├── commands/     # Command handlers (CQRS)
│   ├── queries/      # Query handlers (CQRS)
│   └── ports/        # Interfaces for adapters
└── infrastructure/   # External concerns
    ├── adapters/     # Implementations of ports
    ├── api/         # HTTP handlers and routes
    ├── database/    # Database models and mappers
    └── websocket/   # Real-time communication
```

### Key Components
- **CQRS Implementation**: Separate command and query handlers
- **Repository Pattern**: Abstracted data access with PostgreSQL and in-memory implementations
- **Dependency Injection**: Manual DI in main.go with clear separation of concerns
- **WebSocket Support**: Real-time ledger activity notifications via Gorilla WebSocket
- **HTTP API**: Fiber v2 web framework with CORS enabled

### Entry Points
- **API Server**: `cmd/api/main.go` - Main HTTP server with dependency wiring
- **CLI Tool**: `cmd/cli/cli.go` - Command-line interface (legacy, not currently used)

### Core Domain
LeNavire is a shared expense tracking system supporting:
- **Expenses**: Track shared costs with automatic ID generation and timestamping
- **Payments**: Record payments with contributor information and payment types
- **Ledger Queries**: Real-time balance calculations and contributor rankings
- **Activity Streaming**: WebSocket notifications for ledger changes

### Database Strategy
- **GORM ORM**: For PostgreSQL interactions with auto-migration
- **Dual Storage**: JSON files for development, PostgreSQL for production
- **Repository Abstractions**: Swappable storage implementations (in-memory, file, PostgreSQL)

### Testing Approach
- **BDD Style**: Ginkgo/Gomega for behavior-driven testing
- **Test Doubles**: Stub implementations for all external dependencies
- **Unit Tests**: Focus on command handlers with isolated testing