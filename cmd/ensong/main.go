package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tombell/ensong/internal/config"
	"github.com/tombell/ensong/internal/monitor"
)

func main() {
	home, _ := os.UserHomeDir()
	configPath := flag.String("config", home+"/.config/ensong/ensong.toml", "path to the configuration file")

	logger := log.New(os.Stderr, "", log.LstdFlags)

	flag.Parse()

	cfg, err := config.Load(*configPath)
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
