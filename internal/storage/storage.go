package storage

type Database interface {
	CreateWallet() (string, error)
	GetBalance(walletID string) (int64, error)
	ExecuteTransaction(walletID string, operation string, amount int64) (int64, error)
	Close() error
}
