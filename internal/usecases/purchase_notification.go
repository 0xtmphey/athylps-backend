package usecases

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/biter777/countries"
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

var mapStoreNames = map[string]string{
	"APP_STORE":   "App Store",
	"PLAY_STORE":  "Google Play",
	"STRIPE":      "Stripe",
	"PROMOTIONAL": "RC Manual",
	"RU_STORE":    "RuStore",
}

type SendPurchaseNotificationParams struct {
	EventType     string
	Store         string
	CountryCode   *string
	Price         *float32
	ProductID     *string
	RenewalNumber *int
}

type tgNotifier interface {
	Notify(ctx context.Context, message string) error
}

type SendPurchaseNotificationUsecase struct {
	notifier tgNotifier
	logger   *zap.Logger
}

func NewSendPurchaseNotificationUsecase(notifier tgNotifier, logger *zap.Logger) *SendPurchaseNotificationUsecase {
	return &SendPurchaseNotificationUsecase{
		notifier: notifier,
		logger:   logger,
	}
}

func (u *SendPurchaseNotificationUsecase) Perform(ctx context.Context, params *SendPurchaseNotificationParams) {
	u.logger.Info("Sending purchase notification")

	if !slices.Contains(supportedEventTypes, params.EventType) {
		u.logger.Info("ignoring event type", zap.String("event_type", params.EventType))
		return
	}

	msg := buildNotificationMessage(params)
	err := u.notifier.Notify(ctx, msg)
	if err != nil {
		u.logger.Error("failed to perform notify", zap.Error(err))
	}
}

func buildNotificationMessage(p *SendPurchaseNotificationParams) string {
	var sb strings.Builder

	storeName, ok := mapStoreNames[p.Store]
	if !ok {
		storeName = p.Store
	}

	switch p.EventType {
	case typeInitialPurchase:
		sb.WriteString(fmt.Sprintf("💵 Совершена покупка в <b>%s</b> 💵\n\n", storeName))
	case typeNonRenewingPurchase:
		sb.WriteString(fmt.Sprintf("💵 Совершена разовая покупка в <b>%s</b> 💵\n\n", storeName))
	case typeRenewal:
		sb.WriteString(fmt.Sprintf("🔁 Подписка продлена в <b>%s</b> 🔁\n\n", storeName))
	case typeCancellation:
		sb.WriteString(fmt.Sprintf("✖︎ Совершена отмена подписки в <b>%s</b> ✖︎\n\n", storeName))
	default:
		sb.WriteString(fmt.Sprintf("Произошло событие: %s", p.EventType))
		return sb.String()
	}

	if p.Price != nil && p.EventType != typeCancellation {
		sb.WriteString(fmt.Sprintf("Стоимость: $%.2f\n", *p.Price))
	}

	if p.CountryCode != nil {
		sb.WriteString(fmt.Sprintf("Страна: %s\n", countryName(*p.CountryCode)))
	}

	if p.ProductID != nil {
		sb.WriteString(fmt.Sprintf("Продукт: %s\n", *p.ProductID))
	}

	if p.RenewalNumber != nil {
		sb.WriteString(fmt.Sprintf("Кол-во продлений: %d\n", *p.RenewalNumber))
	}

	return sb.String()
}

func countryName(countryCode string) string {
	country := countries.ByName(countryCode)
	name := country.StringRus()
	flag := country.Emoji()

	return fmt.Sprintf("%s %s", flag, name)
}
