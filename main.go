package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"roulette/internal/config"
	"roulette/internal/lru"
	"roulette/internal/storage"
	"roulette/pkg/server"
)

func main() {
	var configFileLocation string
	flag.StringVar(&configFileLocation, "c", "", "Pass the path for the Roulette config file")
	flag.Parse()

	if configFileLocation == "" {
		fmt.Fprintf(os.Stderr, "You must provide a config file as input")
		os.Exit(-1)
	}

	cfg, err := config.New(configFileLocation)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %s", err)
		os.Exit(-1)
	}

	log.Printf("Config is %v", cfg)
	cache := lru.NewInMemory(float64(cfg.CacheSize))
	storage, err := storage.NewGCS(cfg.Bucket)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialise the storage GCS instance: %s", err)
		os.Exit(-1)
	}

	svr := server.New(cfg, cache, storage)
	err = svr.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialise server: %s", err)
		os.Exit(-1)
	}

}
