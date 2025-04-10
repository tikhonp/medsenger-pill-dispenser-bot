package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/config"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

type command string

const (
	// PrintDbString prints the database configurtion string
	// for sql connection
	PrintDbString command = "print-db-string"

	// AddPillDispenser adds new pill dispenser with serial number and hardware type
	AddPillDispenser command = "add-pill-dispenser"
)

func (c *command) Set(value string) error {
	switch command(value) {
	case PrintDbString, AddPillDispenser:
		*c = command(value)
		return nil
	default:
		return fmt.Errorf("invalid command %s", value)
	}
}

func (c *command) String() string {
	return string(*c)
}

type manageConfig struct {
	command    command
	configPath string

	serialNumber string
	hwType       models.HardwareType
}

func parseFlags() *manageConfig {
	cfg := &manageConfig{}

	const commandUsage = "command to run. Available commands: print-db-string, create-super-admin"
	flag.Var(&cfg.command, "command", commandUsage)
	flag.Var(&cfg.command, "c", commandUsage+" (shorthand)")

	flag.StringVar(&cfg.configPath, "config", "config.pkl", "path to config file")

	flag.StringVar(&cfg.serialNumber, "serial-number", "", "serial number of pill dispenser")
	flag.Var(&cfg.hwType, "hardware-type", "pill dispenser hardware type")

	flag.Parse()

	return cfg
}

func printDbString(cfg *config.Config) {
	fmt.Print(
		cfg.Db.DbFilePath,
	)
}

func addPillDispenser(cfg *config.Config, serialNumber string, hwType models.HardwareType) {
	util.Assert(serialNumber != "", "provide serial number")
	util.Assert(hwType != "", "provide hardware type")

	modelsFactory, err := db.Connect(cfg.Db)
	util.AssertNoErr(err)

	util.AssertNoErr(
		modelsFactory.PillDispensers().New(serialNumber, hwType),
	)
}

func main() {
	manageConfig := parseFlags()

	cfg, err := config.LoadFromPath(context.Background(), manageConfig.configPath)
	util.AssertNoErr(err)

	switch manageConfig.command {
	case PrintDbString:
		printDbString(cfg)
	case AddPillDispenser:
		addPillDispenser(cfg, manageConfig.serialNumber, manageConfig.hwType)
	default:
		fmt.Println("Incvalid arguments")
	}
}
