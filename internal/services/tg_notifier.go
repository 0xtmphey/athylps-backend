package services

import (
	"context"
	"fmt"

	"athylps/internal/config"

	tgbot "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type TgNotifierService struct {
	bot          *tgbot.Bot
	notifyChatID string
	logger       *zap.Logger
}

func NewTgNotifierService(
	cfg *config.TelegramConfig,
	logger *zap.Logger,
) *TgNotifierService {
	bot, err := tgbot.New(cfg.BotToken)
	if err != nil {
		logger.Fatal("failed to create telegram bot", zap.Error(err))
	}
	return &TgNotifierService{
		bot:          bot,
		notifyChatID: cfg.NotifyChatID,
		logger:       logger,
	}
}

func (s *TgNotifierService) Notify(ctx context.Context, message string) error {
	msg, err := s.bot.SendMessage(ctx, &tgbot.SendMessageParams{
		Text:      message,
		ChatID:    s.notifyChatID,
		ParseMode: tgmodels.ParseModeHTML,
	})
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	s.logger.Info("sent telegram message", zap.Any("msg", msg), zap.String("chat_id", s.notifyChatID))

	return nil
}
