package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"wallet-app/internal/application"
	"wallet-app/internal/storage"
)

func TestCreateWallet(t *testing.T) {
	db := storage.NewMockPostgresDB()
	server := httptest.NewServer(application.NewHttpServer(db))
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/v1/wallets", "application/json", nil)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(t, result["walletId"])
}

func TestGetBalance(t *testing.T) {
	db := storage.NewMockPostgresDB()
	router := application.NewHttpServer(db)

	// Create a wallet
	walletID, _ := db.CreateWallet()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result map[string]int64
	json.NewDecoder(rec.Body).Decode(&result)
	assert.Equal(t, int64(0), result["balance"])
}

func TestHandleTransaction(t *testing.T) {
	db := storage.NewMockPostgresDB()
	router := application.NewHttpServer(db)

	walletID, _ := db.CreateWallet()
	body := map[string]interface{}{
		"walletId":      walletID,
		"operationType": "DEPOSIT",
		"amount":        500,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	balance, _ := db.GetBalance(walletID)
	assert.Equal(t, int64(500), balance)
}
