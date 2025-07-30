package main

import (
	"flag"
	"fmt"

	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util/assert"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util/config"
)

type command string

const (
	// PrintDBString prints the database configuration string
	// for sql connection
	PrintDBString command = "print-db-string"

	// AddPillDispenser adds new pill dispenser with serial number and hardware type
	AddPillDispenser command = "add-pill-dispenser"
)

func (c *command) Set(value string) error {
	switch command(value) {
	case PrintDBString, AddPillDispenser:
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

	serialNumber                string
	hwType                      models.HardwareType
	addPillDispenserInteractive bool
}

func parseFlags() *manageConfig {
	cfg := &manageConfig{}

	const commandUsage = "command to run. Available commands: print-db-string, create-super-admin"
	flag.Var(&cfg.command, "command", commandUsage)
	flag.Var(&cfg.command, "c", commandUsage+" (shorthand)")

	flag.StringVar(&cfg.configPath, "config", "config.pkl", "path to config file")

	flag.StringVar(&cfg.serialNumber, "serial-number", "", "serial number of pill dispenser")
	flag.BoolVar(&cfg.addPillDispenserInteractive, "i", false, "prompt for data")
	flag.Var(&cfg.hwType, "hardware-type", "pill dispenser hardware type")

	flag.Parse()

	return cfg
}

func printDBString(cfg *config.Config) {
	fmt.Print(
		db.DataSourceName(cfg.DB),
	)
}

func addPillDispenserInteractive(cfg *config.Config) {
	var serialNumber string
	var hwType models.HardwareType

	fmt.Print("Serial Number: ")
	_, err := fmt.Scanln(&serialNumber)
	assert.NoErr(err)

	fmt.Print("Hardware Type: ")
	_, err = fmt.Scanln(&hwType)
	assert.NoErr(err)

	addPillDispenser(cfg, serialNumber, hwType)
}

func addPillDispenser(cfg *config.Config, serialNumber string, hwType models.HardwareType) {
	assert.C(serialNumber != "", "provide serial number")
	assert.C(hwType == models.HardwareType2x2 || hwType == models.HardwareType4x7, "provide hardware type")

	modelsFactory, err := db.Connect(cfg.DB)
	assert.NoErr(err)

	assert.NoErr(
		modelsFactory.PillDispensers().New(serialNumber, hwType),
	)
}

func main() {
	manageConfig := parseFlags()
	cfg := config.LoadConfigFromEnv()

	switch manageConfig.command {
	case PrintDBString:
		printDBString(cfg)
	case AddPillDispenser:
		if manageConfig.addPillDispenserInteractive {
			addPillDispenserInteractive(cfg)
		} else {
			addPillDispenser(cfg, manageConfig.serialNumber, manageConfig.hwType)
		}
	default:
		fmt.Println("Invalid arguments")
	}
}
