package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/config"
)

type command string

const (
	// PrintDbString prints the database configurtion string
	// for sql connection
	PrintDbString command = "print-db-string"
)

func (c *command) Set(value string) error {
	switch command(value) {
	case PrintDbString:
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
}

func parseFlags() *manageConfig {
	cfg := &manageConfig{}

	const commandUsage = "command to run. Available commands: print-db-string, create-super-admin"
	flag.Var(&cfg.command, "command", commandUsage)
	flag.Var(&cfg.command, "c", commandUsage+" (shorthand)")

	flag.StringVar(&cfg.configPath, "config", "config.pkl", "path to config file")

	flag.Parse()

	return cfg
}

func printDbString(cfg *config.Config) {
	fmt.Print(
		cfg.Db.DbFilePath,
	)
}

func main() {
	manageConfig := parseFlags()

	cfg, err := config.LoadFromPath(context.Background(), manageConfig.configPath)
	if err != nil {
		panic(err)
	}

	switch manageConfig.command {
	case PrintDbString:
		printDbString(cfg)
	}
}
