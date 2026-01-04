package services

import (
	"fmt"
)

type QRCodeService struct {
	verifierEndpoint string
}

func NewQRCodeService(verifierEndpoint string) *QRCodeService {
	return &QRCodeService{
		verifierEndpoint: verifierEndpoint,
	}
}

type QRData struct {
	Type            string `json:"type"`
	ProofRequestID  string `json:"proof_request_id"`
	CallbackURL     string `json:"callback_url"`
	VerifierEndpoint string `json:"verifier_endpoint,omitempty"`
}

func (q *QRCodeService) GenerateQRData(proofRequestID, callbackURL string) (string, error) {
	qrURL := fmt.Sprintf("%s?pres_ex_id=%s&response_uri=%s", 
		q.verifierEndpoint, 
		proofRequestID,
		callbackURL)
	
	return qrURL, nil
}

