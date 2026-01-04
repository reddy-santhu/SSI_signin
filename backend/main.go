package main

import (
	"log"

	"ssi-signin/backend/config"
	"ssi-signin/backend/handlers"
	authMiddleware "ssi-signin/backend/middleware"
	"ssi-signin/backend/repositories"
	"ssi-signin/backend/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	db, err := services.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ariesService := services.NewAriesService(cfg.IssuerAgentURL, cfg.VerifierAgentURL, cfg.LedgerURL)
	verifierService := services.NewVerifierService(cfg.VerifierAgentURL)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	setupRoutes(e, db, ariesService, verifierService)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

func setupRoutes(e *echo.Echo, db *services.Database, ariesService *services.AriesService, verifierService *services.VerifierService) {
	healthHandler := handlers.NewHealthHandler(db, ariesService, verifierService)
	authHandler := handlers.NewAuthHandlerWithDeps(db, ariesService, verifierService)
	credentialHandler := handlers.NewCredentialHandler(ariesService)

	e.GET("/health", healthHandler.Check)

	api := e.Group("/api")
	api.GET("/status", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "SSI Sign-In API",
		})
	})

	api.POST("/login", authHandler.Login)
	api.POST("/proof-callback", authHandler.ProofCallback)
	api.GET("/login/status/:proofRequestId", authHandler.LoginStatus)

	authMW := authMiddleware.NewAuthMiddleware(repositories.NewSessionRepository(db.DB))
	protected := api.Group("", authMW.RequireAuth)
	protected.GET("/dashboard", authHandler.Dashboard)

	api.POST("/schemas", credentialHandler.CreateSchema)
	api.POST("/credential-definitions", credentialHandler.CreateCredentialDefinition)
}

