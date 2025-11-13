package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"athylps/internal/config"
	"athylps/internal/handlers/hooks"
	"athylps/internal/services"
	"athylps/internal/usecases"

	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func Run(
	cfg *config.Config,
	logger *zap.Logger,
	dbpool *pgxpool.Pool,
) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		logger.Fatal("failed to initialize firebase app")
	}

	logger.Info("Initialized Firebase", zap.Any("app", app))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	tgNotifierService := services.NewTgNotifierService(&cfg.Telegram, logger)
	purchaseNotificationUsecase := usecases.NewSendPurchaseNotificationUsecase(tgNotifierService, logger)

	r.Post("/hooks/revenuecat", hooks.HandleRevenueCatWebHook(&cfg.RevenueCat, logger, purchaseNotificationUsecase))
	r.Post("/hooks/rustore", hooks.HandleRustoreWebHook(&cfg.Rustore, logger, purchaseNotificationUsecase))
	r.Get("/hooks/donationalerts", hooks.HandleDonationAlertsWebhook())

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s (environment: %s)", addr, cfg.Server.Env)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
