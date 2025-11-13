package hooks

import (
	"net/http"

	"athylps/internal/config"

	"go.uber.org/zap"
)

func HandleRustoreWebHook(
	cfg *config.RustoreConfig,
	logger *zap.Logger,
	usecase sendPurchaseNotificationUsecase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}
