package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 3000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stage|prod)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	addr := fmt.Sprintf(":%d", cfg.port)

	store, err := NewPostgresStore("postgres", "postgres", "root")
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewApiServer(cfg, addr, logger, store)
	server.Start()
}
