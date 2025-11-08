package hooks

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"athylps/internal/api"
	"athylps/internal/config"

	"go.uber.org/zap"
)

var ErrUnauthorized = errors.New("Unauthorized")

func HandleRevenueCatWebHook(cfg *config.RevenueCatConfig, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		err := validateToken(authHeader, cfg.BearerToken)
		if err != nil {
			logger.Warn("validate token failed", zap.Error(err), zap.String("auth_header", authHeader))
			resp := api.WebhookResponse{Status: api.Success}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		var data api.RevenueCatWebhookEvent
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Warn("failed to decode revenuecat hook request")
			resp := api.WebhookResponse{Status: api.Success}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		logger.Info("[RC webhook]", zap.Any("data", data))
		resp := api.WebhookResponse{Status: api.Success}
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func validateToken(authHeader string, token string) error {
	if authHeader == "" {
		return ErrUnauthorized
	}

	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 {
		return ErrUnauthorized
	}

	headerToken := strings.TrimSpace(splitToken[1])
	if headerToken != token {
		return ErrUnauthorized
	}

	return nil
}
