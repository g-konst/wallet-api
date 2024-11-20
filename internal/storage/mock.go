package storage

import (
	"errors"
	"sync"
)

type MockPostgresDB struct {
	wallets map[string]int64
	mu      sync.Mutex
}

var _ Database = (*MockPostgresDB)(nil)

func NewMockPostgresDB() *MockPostgresDB {
	return &MockPostgresDB{
		wallets: make(map[string]int64),
	}
}

func (m *MockPostgresDB) CreateWallet() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	walletID := "mock-wallet-id"
	m.wallets[walletID] = 0
	return walletID, nil
}

func (m *MockPostgresDB) GetBalance(walletID string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	balance, ok := m.wallets[walletID]
	if !ok {
		return 0, errors.New("wallet not found")
	}
	return balance, nil
}

func (m *MockPostgresDB) ExecuteTransaction(walletID string, operationType string, amount int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	balance, ok := m.wallets[walletID]
	if !ok {
		return balance, errors.New("wallet not found")
	}

	if operationType == "DEPOSIT" {
		m.wallets[walletID] += amount
	} else if operationType == "WITHDRAW" {
		if balance < amount {
			return balance, errors.New("insufficient funds")
		}
		m.wallets[walletID] -= amount
	} else {
		return balance, errors.New("invalid operation type")
	}
	return balance, nil
}

func (m *MockPostgresDB) Close() error {
	return nil
}
