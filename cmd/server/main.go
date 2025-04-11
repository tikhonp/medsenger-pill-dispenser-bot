package main

import (
	"context"

	"github.com/TikhonP/maigo"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/config"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/router"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func main() {
	// Read the configuration from the pkl file
	cfg, err := config.LoadFromPath(context.Background(), "config.pkl")
	util.AssertNoErr(err)

	// Connect to the database
	modelsFactory, err := db.Connect(cfg.Db)
	util.AssertNoErr(err)

	maigoClient := maigo.Init(cfg.Server.MedsengerAgentKey)

	// Setup server
	r := router.New(cfg)
	router.RegisterRoutes(r, cfg, modelsFactory, maigoClient)
	r.Logger.Fatal(router.Start(r, cfg))
}
