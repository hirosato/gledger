# Gledger Architecture

## Overview

Gledger follows **Hexagonal Architecture** (Ports and Adapters) to ensure clean separation of concerns and enable multiple interfaces (CLI, Web UI, etc.) to share the same business logic.

## Directory Structure

```
gledger/
├── domain/              # Core business logic (no external dependencies)
│   ├── ports/          # Interfaces for external dependencies
│   │   ├── parser.go   # Parser port interface
│   │   ├── formatter.go # Formatter port interface
│   │   └── storage.go  # Storage port interface
│   ├── account.go      # Account entity
│   ├── transaction.go  # Transaction entity
│   ├── balance.go      # Balance value object
│   ├── amount.go       # Amount value object
│   └── repository.go   # Repository interfaces
│
├── application/        # Use cases and orchestration
│   ├── dto/           # Data Transfer Objects
│   │   ├── balance_report.go
│   │   ├── account_list.go
│   │   └── register_entry.go
│   ├── usecases/      # Business use cases
│   │   ├── get_balance.go
│   │   ├── list_accounts.go
│   │   └── show_register.go
│   └── journal.go     # Journal aggregate
│
├── adapters/          # Input/Output adapters
│   ├── inbound/       # Input adapters
│   │   └── cli/       # CLI adapter
│   │       ├── commands/     # CLI commands
│   │       └── presenters/   # Output formatting
│   └── outbound/      # Output adapters
│       └── filesystem/       # File system adapter
│           └── parser_adapter.go
│
├── infrastructure/    # Technical implementations
│   └── parser/        # Ledger file parser
│
├── cmd/              # Application entry points
│   └── gledger/      # Main CLI application
│       └── main.go   # Dependency injection setup
│
└── interfaces/       # (Legacy - to be removed)
```

## Key Architectural Principles

### 1. **Dependency Rule**
- Dependencies only point inward
- Domain layer has no external dependencies
- Application layer depends only on domain
- Adapters depend on application and domain
- Infrastructure implements domain ports

### 2. **Port/Adapter Pattern**
- **Ports** (interfaces) are defined in the domain layer
- **Adapters** implement these ports in the infrastructure/adapters layers
- This allows swapping implementations without changing business logic

### 3. **Dependency Injection**
- All dependencies are injected through constructors
- Main.go acts as the composition root
- No hard-coded dependencies between layers

### 4. **Data Transfer Objects (DTOs)**
- Used for communication between layers
- Prevent domain models from leaking to external layers
- Enable different representations for different interfaces

## Data Flow

### CLI Request Flow:
1. **CLI Command** (adapters/inbound/cli) receives user input
2. **Use Case** (application/usecases) orchestrates business logic
3. **Domain Services** execute business rules
4. **Port Implementation** (adapters/outbound) handles I/O
5. **Presenter** (adapters/inbound/cli/presenters) formats output
6. **CLI Response** displayed to user

### Future Web UI Flow:
1. **HTTP Handler** (adapters/inbound/web) receives request
2. **Same Use Case** (application/usecases) orchestrates logic
3. **Same Domain Services** execute rules
4. **Same Port Implementation** handles I/O
5. **JSON Serializer** formats response
6. **HTTP Response** sent to browser

## Benefits

### 1. **Testability**
- Domain logic can be tested in isolation
- Use cases can be tested with mock adapters
- No need for file system or external dependencies in tests

### 2. **Maintainability**
- Clear separation of concerns
- Business logic isolated from technical details
- Easy to understand and modify

### 3. **Extensibility**
- Easy to add new interfaces (Web UI, API, Mobile)
- Can swap implementations (e.g., database instead of files)
- New features don't affect existing code

### 4. **Reusability**
- Same business logic for all interfaces
- DTOs ensure consistent data structures
- Use cases are interface-agnostic

## Adding a Web UI

When adding a web UI in the future:

1. Create `adapters/inbound/web/` directory
2. Add HTTP handlers that call existing use cases
3. Create JSON presenters for DTOs
4. Reuse all application and domain logic
5. No changes needed to business logic

Example:
```go
// adapters/inbound/web/handlers/balance.go
func BalanceHandler(useCase *usecases.GetBalance) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        options := parseOptions(r)
        report, err := useCase.Execute(options)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        json.NewEncoder(w).Encode(report)
    }
}
```

## Future Improvements

1. **Complete Repository Pattern**: Implement repository interfaces in infrastructure
2. **Configuration Management**: Add configuration layer for different environments
3. **Error Handling**: Implement domain-specific error types
4. **Logging**: Add structured logging with proper abstraction
5. **Event Sourcing**: Consider for audit trail and undo/redo functionality