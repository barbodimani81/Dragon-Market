package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/barbodimani81/Dragon-Market/pkg/database"
	"log"
	"net/http"
	"os"

	"github.com/barbodimani81/Dragon-Market/internal/auction/service"
	gen "github.com/barbodimani81/Dragon-Market/internal/database/gen"
	wservice "github.com/barbodimani81/Dragon-Market/internal/wallet/service"
	_ "github.com/lib/pq"
)

type SimpleTransactor struct {
	db *sql.DB
}

// WithinTransaction now strictly implements database.Transactor
func (t *SimpleTransactor) WithinTransaction(ctx context.Context, fn func(database.TransactionQuerier) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Wrap the SQL transaction into your generated sqlc/database queries
	queries := gen.New(tx)

	// Execute the unit of work
	if err := fn(queries); err != nil {
		return err
	}

	return tx.Commit()
}

func main() {
	// 1. Resolve Database connection properties from the environment
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPass := os.Getenv("APP_DATABASE__PASSWORD")
	if dbPass == "" {
		log.Fatalf("Environment variable APP_DATABASE__PASSWORD is required but was not set")
	}

	connStr := fmt.Sprintf("postgres://postgres:%s@%s:5432/dragon_market?sslmode=disable", dbPass, dbHost)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL database.")

	// 2. Initialize Repositories and Domain Services
	queries := gen.New(db)
	transactor := &SimpleTransactor{db: db}

	_ = wservice.NewWalletService(queries)
	_ = service.NewAuctionService(queries, transactor)

	// 3. Start Server Container
	mux := http.NewServeMux()

	// Health check endpoint for verifying connectivity
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	log.Println("Dragon-Market API Server booting up on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
