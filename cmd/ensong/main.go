package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tombell/ensong/pkg/config"
	"github.com/tombell/ensong/pkg/monitor"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	cfg, err := config.Load("./ensong.toml")
	if err != nil {
		logger.Fatalf("error: %s", err)
	}

	m, err := monitor.New(cfg, logger)
	if err != nil {
		logger.Fatalf("error: %s", err)
	}

	go m.Run()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	log.Println("shutting down...")
}
