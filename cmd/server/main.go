package main

import (
	"context"
	"log"
	"os"

	"github.com/TikhonP/maigo"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/config"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/router"
)

func main() {
	// Read the configuration from the pkl file
	cfg, err := config.LoadFromPath(context.Background(), "config.pkl")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Connect to the database
	modelsFactory, err := db.Connect(cfg.Db)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	maigoClient := maigo.Init(cfg.Server.MedsengerAgentKey)

	// Setup server
	r := router.New(cfg)
	router.RegisterRoutes(r, cfg, modelsFactory, maigoClient)
	r.Logger.Fatal(router.Start(r, cfg))
}
