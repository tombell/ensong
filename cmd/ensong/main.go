package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tombell/ensong/pkg/config"
	"github.com/tombell/ensong/pkg/monitor"
)

var configPath = flag.String("config", "ensong.toml", "path to the configuration file")

func main() {
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
