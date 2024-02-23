package main

import (
	"context"
	"findigitalservice/config"
	"findigitalservice/http/rest"
	"log"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("parse .env error")
	}
	server, err := rest.NewServer(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	server.Start(ctx)
}
