package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type VerifierService struct {
	verifierURL string
	client      *http.Client
}

func NewVerifierService(verifierURL string) *VerifierService {
	return &VerifierService{
		verifierURL: verifierURL,
		client:      &http.Client{},
	}
}

func (v *VerifierService) VerifierURL() string {
	return v.verifierURL
}

type ProofRequest struct {
	Name              string                 `json:"name"`
	Version           string                 `json:"version"`
	RequestedAttributes map[string]interface{} `json:"requested_attributes"`
	RequestedPredicates map[string]interface{} `json:"requested_predicates,omitempty"`
}

func (v *VerifierService) CreateProofRequest(proofReq ProofRequest, responseURI string) (string, error) {
	url := fmt.Sprintf("%s/present-proof/create-request", v.verifierURL)
	
	requestedPredicates := proofReq.RequestedPredicates
	if requestedPredicates == nil {
		requestedPredicates = map[string]interface{}{}
	}
	
	payload := map[string]interface{}{
		"proof_request": map[string]interface{}{
			"name":                proofReq.Name,
			"version":             proofReq.Version,
			"requested_attributes": proofReq.RequestedAttributes,
			"requested_predicates": requestedPredicates,
		},
	}
	
	if responseURI != "" {
		payload["response_uri"] = responseURI
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal proof request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := v.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("proof request creation failed: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	proofRequestID, ok := result["presentation_exchange_id"].(string)
	if !ok {
		proofRequestID, ok = result["pres_ex_id"].(string)
	}
	if !ok {
		return "", fmt.Errorf("presentation_exchange_id not found in response: %v", result)
	}

	return proofRequestID, nil
}

func (v *VerifierService) CreateProofRequestWithOOB(proofReq ProofRequest, responseURI string) (string, string, error) {
	proofRequestID, err := v.CreateProofRequest(proofReq, responseURI)
	if err != nil {
		return "", "", err
	}

	proofReqURL := fmt.Sprintf("%s/present-proof/records/%s", v.verifierURL, proofRequestID)
	req, err := http.NewRequest("GET", proofReqURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to get proof request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to get proof request: %s", string(body))
	}

	var proofReqData map[string]interface{}
	if err := json.Unmarshal(body, &proofReqData); err != nil {
		return "", "", fmt.Errorf("failed to parse proof request: %w", err)
	}

	presentationRequestDict, ok := proofReqData["presentation_request_dict"].(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("presentation_request_dict not found in proof request")
	}

	var requestPresentationsAttach []interface{}
	if attachRaw, exists := presentationRequestDict["request_presentations~attach"]; exists {
		switch v := attachRaw.(type) {
		case []interface{}:
			requestPresentationsAttach = v
		case []map[string]interface{}:
			requestPresentationsAttach = make([]interface{}, len(v))
			for i, m := range v {
				requestPresentationsAttach[i] = m
			}
		default:
			return "", "", fmt.Errorf("request_presentations~attach has unexpected type: %T", v)
		}
	} else {
		return "", "", fmt.Errorf("request_presentations~attach not found in presentation_request_dict")
	}
	
	if len(requestPresentationsAttach) == 0 {
		return "", "", fmt.Errorf("request_presentations~attach is empty")
	}

	oobURL := fmt.Sprintf("%s/out-of-band/create-invitation", v.verifierURL)
	
	attachmentType, _ := presentationRequestDict["@type"].(string)
	if attachmentType == "" {
		attachmentType = "did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/present-proof/1.0/request-presentation"
	}
	
	attachments := []map[string]interface{}{
		{
			"@id": "request-0",
			"type": "present-proof",
			"data": map[string]interface{}{
				"json": map[string]interface{}{
					"@type":                     attachmentType,
					"@id":                       presentationRequestDict["@id"],
					"request_presentations~attach": requestPresentationsAttach,
				},
			},
		},
	}

	payload := map[string]interface{}{
		"auto_accept":        true,
		"public":             true,
		"handshake_protocols": []string{},
		"attachments":        attachments,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal OOB invitation: %w", err)
	}

	req, err = http.NewRequest("POST", oobURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create OOB request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err = v.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send OOB request: %w", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read OOB response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		errorMsg := string(body)
		if len(errorMsg) > 200 {
			errorMsg = errorMsg[:200] + "..."
		}
		return "", "", fmt.Errorf("OOB invitation creation failed (status %d): %s. Payload was: %s", resp.StatusCode, errorMsg, string(jsonData))
	}

	var oobResult map[string]interface{}
	if err := json.Unmarshal(body, &oobResult); err != nil {
		return "", "", fmt.Errorf("failed to parse OOB response: %w", err)
	}

	invitationURL, ok := oobResult["invitation_url"].(string)
	if !ok {
		return "", "", fmt.Errorf("invitation_url not found in OOB response: %v", oobResult)
	}

	verifierEndpoint := os.Getenv("VERIFIER_ENDPOINT")
	if verifierEndpoint != "" {
		if strings.Contains(invitationURL, "verifier-agent:8003") {
			invitationURL = strings.ReplaceAll(invitationURL, "verifier-agent:8003", verifierEndpoint)
		}
		if strings.Contains(invitationURL, "localhost:8003") {
			invitationURL = strings.ReplaceAll(invitationURL, "localhost:8003", verifierEndpoint)
		}
	}

	return proofRequestID, invitationURL, nil
}

func (v *VerifierService) VerifyProof(proofRequestID string, proof map[string]interface{}) (bool, error) {
	url := fmt.Sprintf("%s/present-proof/records/%s/verify-presentation", v.verifierURL, proofRequestID)
	
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := v.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("proof verification failed: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("failed to parse response: %w", err)
	}

	verified, ok := result["verified"].(bool)
	if !ok {
		state, _ := result["state"].(string)
		return state == "verified", nil
	}

	return verified, nil
}

