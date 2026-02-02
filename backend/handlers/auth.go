package handlers

import (
	"crypto/subtle"
	"log"
	"net/http"
	"ssi-signin/backend/config"
	"ssi-signin/backend/models"
	"ssi-signin/backend/repositories"
	"ssi-signin/backend/services"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	cfg               *config.Config
	db                *services.Database
	ariesService      *services.AriesService
	verifierService   *services.VerifierService
	qrService         *services.QRCodeService
	sessionService    *services.SessionService
	proofRequestStore *services.ProofRequestStore
	userRepo          *repositories.UserRepository
	sessionRepo       *repositories.SessionRepository
}

func NewAuthHandlerWithDeps(db *services.Database, cfg *config.Config, ariesService *services.AriesService, verifierService *services.VerifierService) *AuthHandler {
	return &AuthHandler{
		cfg:               cfg,
		db:                db,
		ariesService:      ariesService,
		verifierService:   verifierService,
		qrService:         services.NewQRCodeService(cfg.VerifierPublicURL),
		sessionService:    services.NewSessionService(),
		proofRequestStore: services.NewProofRequestStore(),
		userRepo:          repositories.NewUserRepository(db.DB),
		sessionRepo:       repositories.NewSessionRepository(db.DB),
	}
}

type LoginRequest struct {
	CallbackURL string `json:"callback_url"`
}

type LoginResponse struct {
	QRData         string `json:"qr_data"`
	ProofRequestID string `json:"proof_request_id"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	credDefID := h.cfg.CredentialDefinitionID
	if credDefID == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "CREDENTIAL_DEFINITION_ID environment variable not set",
		})
	}

	proofReq := services.ProofRequest{
		Name:    "SSI Sign-In",
		Version: "1.0",
		RequestedAttributes: map[string]interface{}{
			"kyc_verified": map[string]interface{}{
				"name": "kyc_verified",
				"restrictions": []map[string]interface{}{
					{
						"cred_def_id": credDefID,
					},
				},
			},
		},
	}

	callbackURL := c.QueryParam("callback_url")
	if callbackURL == "" {
		callbackURL = h.cfg.CallbackURL
	}

	proofRequestID, invitationURL, err := h.verifierService.CreateProofRequestWithOOB(proofReq, callbackURL)
	if err != nil {
		proofRequestID, err = h.verifierService.CreateProofRequest(proofReq, callbackURL)
		if err != nil {
			log.Printf("login create proof request: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to create proof request",
			})
		}

		h.proofRequestStore.Set(proofRequestID, "")

		qrData, err := h.qrService.GenerateQRData(proofRequestID, callbackURL)
		if err != nil {
			log.Printf("login generate QR: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to generate QR code",
			})
		}

		return c.JSON(http.StatusOK, LoginResponse{
			QRData:         qrData,
			ProofRequestID: proofRequestID,
		})
	}

	h.proofRequestStore.Set(proofRequestID, "")

	qrData := invitationURL

	return c.JSON(http.StatusOK, LoginResponse{
		QRData:         qrData,
		ProofRequestID: proofRequestID,
	})
}

type ProofCallbackRequest struct {
	ProofRequestID string                 `json:"proof_request_id"`
	Proof          map[string]interface{} `json:"proof"`
	HolderDID      string                 `json:"holder_did"`
}

type ProofCallbackResponse struct {
	Success      bool   `json:"success"`
	SessionToken string `json:"session_token,omitempty"`
	Message      string `json:"message"`
}

func (h *AuthHandler) ProofCallback(c echo.Context) error {
	if want := h.cfg.ProofCallbackSecret; want != "" {
		got := c.Request().Header.Get("X-Webhook-Secret")
		if len(got) != len(want) || subtle.ConstantTimeCompare([]byte(got), []byte(want)) != 1 {
			log.Printf("proof-callback: invalid webhook secret")
			return c.JSON(http.StatusUnauthorized, ProofCallbackResponse{
				Success: false,
				Message: "Unauthorized",
			})
		}
	}

	var req ProofCallbackRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("proof-callback bind: %v", err)
		return c.JSON(http.StatusBadRequest, ProofCallbackResponse{
			Success: false,
			Message: "Invalid request",
		})
	}

	verified, err := h.verifierService.VerifyProof(req.ProofRequestID, req.Proof)
	if err != nil {
		log.Printf("proof-callback verify: %v", err)
		return c.JSON(http.StatusInternalServerError, ProofCallbackResponse{
			Success: false,
			Message: "Proof verification failed",
		})
	}

	if !verified {
		return c.JSON(http.StatusUnauthorized, ProofCallbackResponse{
			Success: false,
			Message: "Proof verification failed",
		})
	}

	user, err := h.userRepo.FindByDID(req.HolderDID)
	if err != nil {
		log.Printf("proof-callback find user: %v", err)
		return c.JSON(http.StatusInternalServerError, ProofCallbackResponse{
			Success: false,
			Message: "Database error",
		})
	}

	if user == nil {
		user = &models.User{
			DID: req.HolderDID,
		}
		if err := h.userRepo.Create(user); err != nil {
			log.Printf("proof-callback create user: %v", err)
			return c.JSON(http.StatusInternalServerError, ProofCallbackResponse{
				Success: false,
				Message: "Failed to create user",
			})
		}
	}

	token, err := h.sessionService.GenerateToken()
	if err != nil {
		log.Printf("proof-callback session token: %v", err)
		return c.JSON(http.StatusInternalServerError, ProofCallbackResponse{
			Success: false,
			Message: "Failed to generate session token",
		})
	}

	session := &models.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: h.sessionService.GetExpirationTime(),
	}

	if err := h.sessionRepo.Create(session); err != nil {
		log.Printf("proof-callback create session: %v", err)
		return c.JSON(http.StatusInternalServerError, ProofCallbackResponse{
			Success: false,
			Message: "Failed to create session",
		})
	}

	h.proofRequestStore.Set(req.ProofRequestID, token)

	return c.JSON(http.StatusOK, ProofCallbackResponse{
		Success:      true,
		SessionToken: token,
		Message:      "Login successful",
	})
}

type DashboardResponse struct {
	User *models.User `json:"user"`
}

func (h *AuthHandler) Logout(c echo.Context) error {
	sess, ok := c.Get("session").(*models.Session)
	if !ok || sess == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "session not found",
		})
	}
	if err := h.sessionRepo.DeleteByToken(sess.Token); err != nil {
		log.Printf("logout: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not sign out",
		})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) Dashboard(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user ID from session",
		})
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("dashboard find user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user",
		})
	}

	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	return c.JSON(http.StatusOK, DashboardResponse{
		User: user,
	})
}

type LoginStatusResponse struct {
	Status       string `json:"status"`
	SessionToken string `json:"session_token,omitempty"`
}

func (h *AuthHandler) LoginStatus(c echo.Context) error {
	proofRequestID := c.Param("proofRequestId")

	token, found := h.proofRequestStore.Get(proofRequestID)
	if !found {
		return c.JSON(http.StatusOK, LoginStatusResponse{
			Status: "not_found",
		})
	}

	if token == "" {
		return c.JSON(http.StatusOK, LoginStatusResponse{
			Status: "pending",
		})
	}

	h.proofRequestStore.Delete(proofRequestID)

	return c.JSON(http.StatusOK, LoginStatusResponse{
		Status:       "completed",
		SessionToken: token,
	})
}
