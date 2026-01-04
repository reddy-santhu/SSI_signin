package config

import (
	"os"
)

type Config struct {
	Port             string
	DatabaseURL      string
	IssuerAgentURL   string
	VerifierAgentURL string
	LedgerURL        string
}

func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ssi_db?sslmode=disable"),
		IssuerAgentURL:   getEnv("ISSUER_AGENT_URL", "http://localhost:8002"),
		VerifierAgentURL: getEnv("VERIFIER_AGENT_URL", "http://localhost:8004"),
		LedgerURL:        getEnv("LEDGER_URL", "http://localhost:9000"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
