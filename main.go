package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"wallet-app/internal/handler"
	"wallet-app/internal/service"
	"wallet-app/internal/storage"
)

func main() {
	// Load environment variables
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading env file: %v", err)
	}

	// Initialize database
	db, err := storage.NewPostgresDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize services and handlers
	walletService := service.NewWalletService(db)
	walletHandler := handler.NewWalletHandler(walletService)

	// Set up router
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/wallet", walletHandler.HandleTransaction)
		r.Get("/wallets/{walletId}", walletHandler.GetBalance)
		r.Post("/wallets", walletHandler.CreateWallet) // Новый маршрут
	})

	// Start server
	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
