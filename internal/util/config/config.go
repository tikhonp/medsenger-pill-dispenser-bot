// Package config provides configuration structures and functions to load
package config

import (
	"os"
	"strconv"
)

type Config struct {
	DB *Database `pkl:"db"`

	Server *Server `pkl:"server"`

	Host string `pkl:"host"`

	SentryDSN string `pkl:"sentryDsn"`
}

type Database struct {
	User string

	Password string

	Dbname string

	Host string
}

type Server struct {
	// The port to listen on.
	Port uint16 `pkl:"port"`

	// Sets server to debug mode.
	Debug bool `pkl:"debug"`

	// Medsenger Agent secret key.
	MedsengerAgentKey string `pkl:"medsengerAgentKey"`
}

func LoadConfigFromEnv() *Config {
	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		panic(err)
	}
	return &Config{
		Server: &Server{
			Port:              uint16(serverPort),
			MedsengerAgentKey: os.Getenv("PILLDISPENSER_KEY"),
			Debug:             os.Getenv("DEBUG") == "true",
		},
		DB: &Database{
			User:     os.Getenv("DB_LOGIN"),
			Password: os.Getenv("DB_PASSWORD"),
			Dbname:   os.Getenv("DB_DATABASE"),
			Host:     os.Getenv("DB_HOST"),
		},
		SentryDSN: os.Getenv("SENTRY_DSN"),
	}
}
