# Sequence Diagram: Place Bid Concurrency Flow

Below is the transactional flow showing how the `auction` and `wallet` modules safely interact using transactional closures and row-level database locks.

```mermaid
sequenceDiagram
    autonumber
    actor User as Bidder
    participant API as API Handler
    participant App as Auction Application Service
    participant Tx as Transactor (DB)
    participant RepoW as Wallet Repository
    participant RepoA as Auction Repository

    User->>API: POST /auctions/{id}/bids (Bid Amount)
    API->>App: PlaceBid(ctx, bidderID, auctionID, amount)
    
    App->>Tx: WithinTransaction(ctx, closure)
    Tx->>Tx: BEGIN TRANSACTION
    
    Note over Tx, RepoW: Step 1: Acquire exclusive lock on Wallet
    Tx->>RepoW: GetWalletForUpdate(ctx, bidderID)
    RepoW-->>Tx: Return Wallet (Locked)
    
    Note over Tx, RepoA: Step 2: Acquire exclusive lock on Auction
    Tx->>RepoA: GetAuctionForUpdate(ctx, auctionID)
    RepoA-->>Tx: Return Auction (Locked)

    Note over App: Step 3: Domain Validation & State Transitions
    App->>App: Validate Bid (Amount > Highest Bid? Wallet holds enough balance?)

    Note over Tx, RepoW: Step 4: Perform Atomic Updates
    Tx->>RepoW: ReserveBalance(bidderID, amount)
    Tx->>RepoA: UpdateHighestBid(auctionID, bidderID, amount)

    Tx->>Tx: COMMIT TRANSACTION (Releases all locks)
    Tx-->>App: Success
    App-->>API: Bid Accepted
    API-->>User: 201 Created
