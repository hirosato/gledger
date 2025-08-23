# Domain Layer

This directory contains the core business logic and domain models, completely independent of external dependencies.

## Contents

- **Entities**: Core business objects (Account, Transaction, Posting, Commodity)
- **Value Objects**: Immutable objects (Amount, Balance, Date)
- **Aggregates**: Domain object clusters (Journal)
- **Domain Services**: Business logic that doesn't fit within a single entity
- **Interfaces**: Port definitions for infrastructure dependencies

## Key Principles

- No external dependencies
- Pure business logic
- All business rules and validations
- Domain-driven design patterns