package main

import (
	"github.com/tikhonp/maigo"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/bviews"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/router"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util/assert"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util/config"
)

func initDependencies() util.Dependencies {
	cfg := config.LoadConfigFromEnv()

	if !cfg.Server.Debug {
		err := util.StartSentry(cfg.SentryDSN)
		assert.NoErr(err)
	}

	// Connect to the database
	modelsFactory, err := db.Connect(cfg.DB)
	assert.NoErr(err)

	maigoClient := maigo.Init(cfg.Server.MedsengerAgentKey)

	return util.NewDependencies(cfg, maigoClient, modelsFactory)
}

func main() {
	deps := initDependencies()

	bviews.Host = deps.Cfg.Host

	r := router.New(deps.Cfg)
	router.RegisterRoutes(r, deps)
	r.Logger.Fatal(router.Start(r, deps.Cfg))
}
