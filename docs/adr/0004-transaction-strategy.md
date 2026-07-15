# ADR 0004: Transaction Management across Domain Boundaries

## Status
Proposed

## Context
In our Modular Monolith, we have distinct business domains (e.g., `wallet` and `auction`). A core transaction—such as placing a bid—requires debiting/reserving money in the `wallet` module and registering a bid in the `auction` module.

We must execute these actions inside a single ACID database transaction. However, passing database connection pointers (like `*sql.Tx`) directly into our Domain or Application services violates clean architecture and tightens domain coupling.

## Decision
We will implement a **Unit of Work** or **Transactional Runner** pattern using Go closures.

1. The Application/Service layer will declare a transaction manager interface:
   ```go
   type Transactor interface {
       WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
   }
