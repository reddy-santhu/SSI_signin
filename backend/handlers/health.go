package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"ssi-signin/backend/services"
)

type HealthHandler struct {
	db              *services.Database
	ariesService    *services.AriesService
	verifierService *services.VerifierService
}

func NewHealthHandler(db *services.Database, ariesService *services.AriesService, verifierService *services.VerifierService) *HealthHandler {
	return &HealthHandler{
		db:              db,
		ariesService:    ariesService,
		verifierService: verifierService,
	}
}

func (h *HealthHandler) Check(c echo.Context) error {
	status := map[string]interface{}{
		"status": "ok",
		"checks": make(map[string]string),
	}

	checks := status["checks"].(map[string]string)

	if err := h.db.DB.Ping(); err != nil {
		checks["database"] = "failed: " + err.Error()
		status["status"] = "degraded"
	} else {
		checks["database"] = "ok"
	}

	if err := h.checkAgent(h.ariesService.IssuerURL()); err != nil {
		checks["issuer_agent"] = "failed: " + err.Error()
		status["status"] = "degraded"
	} else {
		checks["issuer_agent"] = "ok"
	}

	if err := h.checkAgent(h.verifierService.VerifierURL()); err != nil {
		checks["verifier_agent"] = "failed: " + err.Error()
		status["status"] = "degraded"
	} else {
		checks["verifier_agent"] = "ok"
	}

	return c.JSON(http.StatusOK, status)
}

func (h *HealthHandler) checkAgent(url string) error {
	resp, err := http.Get(url + "/status")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("agent returned status %d", resp.StatusCode)
	}

	return nil
}

