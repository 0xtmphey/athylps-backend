package app

import (
	"fmt"
	"log"
	"net/http"

	"athylps/internal/config"
	"athylps/internal/handlers/hooks"
	"athylps/internal/services"
	"athylps/internal/usecases"

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
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	tgNotifierService := services.NewTgNotifierService(&cfg.Telegram, logger)
	purchaseNotificationUsecase := usecases.NewSendPurchaseNotificationUsecase(tgNotifierService, logger)

	r.Post("/hooks/revenuecat", hooks.HandleRevenueCatWebHook(&cfg.RevenueCat, logger, purchaseNotificationUsecase))

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s (environment: %s)", addr, cfg.Server.Env)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
