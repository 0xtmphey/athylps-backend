package hooks

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"athylps/internal/api"
	"athylps/internal/config"
	"athylps/internal/usecases"

	"go.uber.org/zap"
)

var ErrUnauthorized = errors.New("Unauthorized")

type sendPurchaseNotificationUsecase interface {
	Perform(params *usecases.SendPurchaseNotificationParams)
}

func HandleRevenueCatWebHook(
	cfg *config.RevenueCatConfig,
	logger *zap.Logger,
	usecase sendPurchaseNotificationUsecase,
) http.HandlerFunc {
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

		usecase.Perform(&usecases.SendPurchaseNotificationParams{
			EventType:     string(data.Event.Type),
			CountryCode:   data.Event.CountryCode,
			Price:         data.Event.Price,
			ProductID:     data.Event.ProductId,
			RenewalNumber: data.Event.RenewalNumber,
			Store:         string(data.Event.Store),
		})

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
