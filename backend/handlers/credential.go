package handlers

import (
	"net/http"
	"ssi-signin/backend/services"

	"github.com/labstack/echo/v4"
)

type CredentialHandler struct {
	ariesService *services.AriesService
}

func NewCredentialHandler(ariesService *services.AriesService) *CredentialHandler {
	return &CredentialHandler{
		ariesService: ariesService,
	}
}

type CreateSchemaRequest struct {
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	Attributes []string `json:"attributes"`
}

type CreateSchemaResponse struct {
	SchemaID string `json:"schema_id"`
	Message  string `json:"message"`
}

func (h *CredentialHandler) CreateSchema(c echo.Context) error {
	var req CreateSchemaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request: " + err.Error(),
		})
	}

	schema := services.Schema{
		Name:       req.Name,
		Version:    req.Version,
		Attributes: req.Attributes,
	}

	schemaID, err := h.ariesService.CreateSchema(schema)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create schema: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, CreateSchemaResponse{
		SchemaID: schemaID,
		Message:  "Schema created successfully",
	})
}

type CreateCredDefRequest struct {
	SchemaID         string `json:"schema_id"`
	SupportRevocation bool  `json:"support_revocation"`
}

type CreateCredDefResponse struct {
	CredDefID string `json:"credential_definition_id"`
	Message   string `json:"message"`
}

func (h *CredentialHandler) CreateCredentialDefinition(c echo.Context) error {
	var req CreateCredDefRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request: " + err.Error(),
		})
	}

	credDefID, err := h.ariesService.CreateCredentialDefinition(req.SchemaID, req.SupportRevocation)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create credential definition: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, CreateCredDefResponse{
		CredDefID: credDefID,
		Message:   "Credential definition created successfully",
	})
}

