package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"wallet-app/internal/handler"
	"wallet-app/internal/service"
	"wallet-app/internal/storage"
)

func NewHttpServer(db storage.Database) http.Handler {
	r := chi.NewRouter()
	setMiddlewares(r)

	walletService := service.NewWalletService(db)
	walletHandler := handler.NewWalletHandler(walletService)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/wallets", walletHandler.CreateWallet)
		r.Get("/wallets/{walletId}", walletHandler.GetBalance)
		r.Post("/wallet", walletHandler.HandleTransaction)
	})

	return r
}

func setMiddlewares(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.Recoverer)

}
