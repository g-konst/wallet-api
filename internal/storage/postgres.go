package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"wallet-app/pkg/logger"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db  *sql.DB
	log *logger.Log
}

var _ Database = (*PostgresDB)(nil)

func NewPostgresDB(dsn string) (*PostgresDB, error) {
	log := logger.NewLogger()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database")
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database")
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS wallets (
			id UUID PRIMARY KEY,
			balance BIGINT NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table")
	}

	return &PostgresDB{db: db, log: log}, nil
}

func (pg *PostgresDB) CreateWallet() (string, error) {
	walletID := uuid.New().String()
	query := "INSERT INTO wallets (id, balance) VALUES ($1, $2)"
	_, err := pg.db.Exec(query, walletID, 0)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet")
	}
	return walletID, nil
}

func (pg *PostgresDB) GetBalance(walletID string) (int64, error) {
	var balance int64
	query := "SELECT balance FROM wallets WHERE id = $1"
	err := pg.db.QueryRow(query, walletID).Scan(&balance)
	return balance, err
}

func (pg *PostgresDB) ExecuteTransaction(walletID string, operation string, amount int64) (int64, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var balance int64
	query := "SELECT balance FROM wallets WHERE id = $1 FOR UPDATE"
	err = tx.QueryRow(query, walletID).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return balance, fmt.Errorf("wallet not found")
		}
		return balance, err
	}

	switch operation {
	case "DEPOSIT":
		balance += amount
	case "WITHDRAW":
		if balance < amount {
			return balance, fmt.Errorf("insufficient funds")
		}
		balance -= amount
	default:
		return balance, fmt.Errorf("invalid operation type")
	}

	query = "UPDATE wallets SET balance = $2 WHERE id = $1"
	_, err = tx.Exec(query, walletID, balance)
	if err != nil {
		return balance, err
	}

	return balance, tx.Commit()
}

func (pg *PostgresDB) Close() error {
	return pg.db.Close()
}
