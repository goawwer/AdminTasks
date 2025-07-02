package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/goawwer/admintasks/config"
	"github.com/goawwer/admintasks/internal/api"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Println("config successfully loaded")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := api.New(config)

	if err := server.Start(ctx); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
