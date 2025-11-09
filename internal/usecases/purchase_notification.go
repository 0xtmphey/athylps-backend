package usecases

import (
	"fmt"
	"slices"
	"strings"

	"go.uber.org/zap"
)

var (
	typeInitialPurchase     = "INITIAL_PURCHASE"
	typeNonRenewingPurchase = "NON_RENEWING_PURCHASE"
	typeRenewal             = "RENEWAL"
	typeCancellation        = "CANCELLATION"
)

var supportedEventTypes = []string{
	typeInitialPurchase,
	typeNonRenewingPurchase,
	typeRenewal,
	typeCancellation,
}

type SendPurchaseNotificationParams struct {
	EventType     string
	Store         string
	CountryCode   *string
	Price         *float32
	ProductID     *string
	RenewalNumber *int
}

type SendPurchaseNotificationUsecase struct {
	logger *zap.Logger
}

func NewSendPurchaseNotificationUsecase(logger *zap.Logger) *SendPurchaseNotificationUsecase {
	return &SendPurchaseNotificationUsecase{logger: logger}
}

func (u *SendPurchaseNotificationUsecase) Perform(params *SendPurchaseNotificationParams) {
	u.logger.Info("Sending purchase notification")

	if !slices.Contains(supportedEventTypes, params.EventType) {
		u.logger.Info("ignoring event type", zap.String("event_type", params.EventType))
		return
	}

	u.logger.Info(buildNotificationMessage(params))
}

func buildNotificationMessage(p *SendPurchaseNotificationParams) string {
	var sb strings.Builder

	switch p.EventType {
	case typeInitialPurchase:
		sb.WriteString("💵 Совершена покупка 💵\n\n")
	case typeNonRenewingPurchase:
		sb.WriteString("💵 Совершена покупка (без продления) 💵\n\n")
	case typeRenewal:
		sb.WriteString("🔁 Подписка продлена 🔁\n\n")
	case typeCancellation:
		sb.WriteString("✖︎ Совершена отмена подписки ✖︎\n\n")
	default:
		sb.WriteString(fmt.Sprintf("Произошло событие: %s", p.EventType))
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("Магазин: %s\n", p.Store))

	if p.Price != nil && p.EventType != typeCancellation {
		sb.WriteString(fmt.Sprintf("Стоимость: $%.2f\n", *p.Price))
	}

	if p.CountryCode != nil {
		sb.WriteString(fmt.Sprintf("Страна: %s\n", *p.CountryCode))
	}

	if p.ProductID != nil {
		sb.WriteString(fmt.Sprintf("Продукт: %s\n", *p.ProductID))
	}

	if p.RenewalNumber != nil {
		sb.WriteString(fmt.Sprintf("Кол-во продлений: %d\n", *p.RenewalNumber))
	}

	return sb.String()
}
