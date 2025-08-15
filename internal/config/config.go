package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

func Load() *Config {
	var cfg Config

	flag.StringVar(&cfg.RunAddress, "a", getEnv("RUN_ADDRESS", ":8080"), "Address and port to run server")
	flag.StringVar(&cfg.DatabaseURI, "d", getEnv("DATABASE_URI", "postgres://localhost:5432/gopher"), "Database connection address")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", getEnv("ACCRUAL_SYSTEM_ADDRESS", ""), "Accrual system address")

	flag.Parse()

	return &cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
