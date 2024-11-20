package main

import (
	"wallet-app/internal"
	"wallet-app/pkg/logger"
)

func main() {
	log := logger.NewLogger()
	cfg := internal.MustLoadConfig()
	if err := internal.Run(cfg, log); err != nil {
		panic(log.WithError(err, "failed to run server"))
	}
}
