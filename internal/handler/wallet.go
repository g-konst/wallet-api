package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"wallet-app/internal/service"
)

type WalletHandler struct {
	service *service.WalletService
}

func NewWalletHandler(s *service.WalletService) *WalletHandler {
	return &WalletHandler{service: s}
}

func (h *WalletHandler) HandleTransaction(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WalletID  string `json:"walletId"`
		Operation string `json:"operationType"`
		Amount    int64  `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "amount must be a positive value", http.StatusBadRequest)
		return
	}

	balance, err := h.service.HandleTransaction(req.WalletID, req.Operation, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int64{"new_balance": balance})
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "walletId")

	balance, err := h.service.GetBalance(walletID)
	if err != nil {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int64{"balance": balance})
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	walletID, err := h.service.CreateWallet()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"walletId": walletID})
}
