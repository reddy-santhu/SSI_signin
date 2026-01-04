package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type AriesService struct {
	issuerURL   string
	verifierURL string
	ledgerURL   string
	client      *http.Client
}

func NewAriesService(issuerURL, verifierURL, ledgerURL string) *AriesService {
	return &AriesService{
		issuerURL:   issuerURL,
		verifierURL: verifierURL,
		ledgerURL:   ledgerURL,
		client:      &http.Client{},
	}
}

func (a *AriesService) IssuerURL() string {
	return a.issuerURL
}

type Schema struct {
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	Attributes []string `json:"attributes"`
}

func (a *AriesService) CreateSchema(schema Schema) (string, error) {
	publicDidURL := fmt.Sprintf("%s/wallet/did/public", a.issuerURL)
	publicDidResp, err := a.client.Get(publicDidURL)
	if err != nil {
		return "", fmt.Errorf("failed to get public DID: %w", err)
	}
	defer publicDidResp.Body.Close()

	var publicDidResult map[string]interface{}
	if err := json.NewDecoder(publicDidResp.Body).Decode(&publicDidResult); err != nil {
		return "", fmt.Errorf("failed to decode public DID response: %w", err)
	}

	var issuerID string
	if result, ok := publicDidResult["result"].(map[string]interface{}); ok {
		did, _ := result["did"].(string)
		if strings.HasPrefix(did, "did:sov:") {
			shortDid := strings.TrimPrefix(did, "did:sov:")
			issuerID = fmt.Sprintf("did:indy:test:%s", shortDid)
		} else if strings.HasPrefix(did, "did:indy:") {
			issuerID = did
		} else if did != "" {
			issuerID = fmt.Sprintf("did:indy:test:%s", did)
		} else {
			return "", fmt.Errorf("no DID found in public DID response")
		}
	}

	if issuerID == "" {
		return "", fmt.Errorf("no public DID found in issuer agent")
	}

	url := fmt.Sprintf("%s/schemas", a.issuerURL)

	payload := map[string]interface{}{
		"schema_name":    schema.Name,
		"schema_version": schema.Version,
		"attributes":     schema.Attributes,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal schema: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("schema creation failed: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	schemaID, ok := result["schema_id"].(string)
	if !ok {
		return "", fmt.Errorf("schema_id not found in response")
	}

	return schemaID, nil
}

func (a *AriesService) CreateCredentialDefinition(schemaID string, supportRevocation bool) (string, error) {
	url := fmt.Sprintf("%s/credential-definitions", a.issuerURL)

	payload := map[string]interface{}{
		"schema_id":          schemaID,
		"support_revocation": supportRevocation,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cred def: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("credential definition creation failed: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	credDefID, ok := result["credential_definition_id"].(string)
	if !ok {
		return "", fmt.Errorf("credential_definition_id not found in response")
	}

	return credDefID, nil
}
