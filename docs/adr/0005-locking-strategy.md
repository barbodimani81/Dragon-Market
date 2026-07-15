This file addresses the **concurrency and race condition** worries directly. It shows you know exactly how to handle high-concurrency contention (like millions of users bidding on a single Legendary item).

```markdown
# ADR 0005: Concurrency and Locking Strategy for High-Contention Operations

## Status
Proposed

## Context
High-value operations—such as placing a bid on an auction or buying a scarce marketplace item—suffer from write-skew and race conditions if multiple concurrent requests read state (e.g., current highest bid) and subsequently write updates based on stale values.

## Decision
We will employ a pessimistic concurrency control strategy using PostgreSQL explicit row locking (`SELECT ... FOR UPDATE`).

1. **Transaction Pipeline:** Every mutating operation on a resource (Auction, Wallet, Marketplace Listing) must follow this strict sequence:
   - Begin Transaction.
   - Acquire an exclusive write lock on the root aggregate immediately using `SELECT ... FOR UPDATE` (e.g., lock the Auction row or the Wallet row by ID).
   - Perform domain validation (Is the auction still open? Is the bid higher than the current highest bid? Does the wallet have enough funds?).
   - Apply state modifications.
   - Commit Transaction (automatically releasing the lock).
   
2. **Deadlock Prevention:** To completely eliminate database deadlocks, locks must **always** be acquired in a strict, uniform system-wide resource order:
   1. `Wallet` rows (lock the bidder/buyer's wallet first).
   2. `Auction` or `Marketplace Listing` rows second.

## Consequences
- **Pros:** Guarantees 100% consistency and prevents any double-spend or invalid bid scenarios at the database level.
- **Cons:** Holds database connections slightly longer under high contention. We will mitigate this by keeping database transactions extremely short and devoid of any external I/O or network calls.
