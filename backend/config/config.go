package config

import (
	"os"
)

type Config struct {
	Port                   string
	DatabaseURL            string
	IssuerAgentURL         string
	VerifierAgentURL       string
	LedgerURL              string
	CredentialDefinitionID string
	CallbackURL            string
	VerifierPublicURL      string
	ProofCallbackSecret    string
}

func Load() *Config {
	return &Config{
		Port:                   getEnv("PORT", "8080"),
		DatabaseURL:            getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ssi_db?sslmode=disable"),
		IssuerAgentURL:         getEnv("ISSUER_AGENT_URL", "http://localhost:8000"),
		VerifierAgentURL:       getEnv("VERIFIER_AGENT_URL", "http://localhost:8002"),
		LedgerURL:              getEnv("LEDGER_URL", "http://localhost:9000"),
		CredentialDefinitionID: getEnv("CREDENTIAL_DEFINITION_ID", ""),
		CallbackURL:            getEnv("CALLBACK_URL", "http://localhost:8080/api/proof-callback"),
		VerifierPublicURL:      getEnv("VERIFIER_ENDPOINT", "http://localhost:8003"),
		ProofCallbackSecret:    getEnv("PROOF_CALLBACK_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
