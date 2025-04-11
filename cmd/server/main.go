package main

import (
	"context"

	"github.com/TikhonP/maigo"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/config"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/router"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func initDependencies() util.Dependencies {
	// Read the configuration from the pkl file
	cfg, err := config.LoadFromPath(context.Background(), "config.pkl")
	util.AssertNoErr(err)

	// Connect to the database
	modelsFactory, err := db.Connect(cfg.Db)
	util.AssertNoErr(err)

	maigoClient := maigo.Init(cfg.Server.MedsengerAgentKey)

	return util.NewDependencies(cfg, maigoClient, modelsFactory)
}

func main() {
	deps := initDependencies()

	// Setup server
	r := router.New(deps.Cfg)
	router.RegisterRoutes(r, deps)
	r.Logger.Fatal(router.Start(r, deps.Cfg))
}
