package service

import (
	"wallet-app/internal/storage"
)

type WalletService struct {
	db storage.Database
}

func NewWalletService(db storage.Database) *WalletService {
	return &WalletService{db: db}
}

func (s *WalletService) HandleTransaction(walletID, operation string, amount int64) (int64, error) {
	return s.db.ExecuteTransaction(walletID, operation, amount)
}

func (s *WalletService) GetBalance(walletID string) (int64, error) {
	return s.db.GetBalance(walletID)
}

func (s *WalletService) CreateWallet() (string, error) {
	return s.db.CreateWallet()
}
