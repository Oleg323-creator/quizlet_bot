package tg

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"quizlet_bot/internal/domain/manager/ucase_manager"
	"sync"
)

type TgBot struct {
	usecases *ucase_manager.ManagerUsecases
	botTg    *tgbotapi.BotAPI
	ctx      context.Context
	wg       *sync.WaitGroup
}

func NewTgBot(usecases *ucase_manager.ManagerUsecases, token string, ctx context.Context, wg *sync.WaitGroup) (*TgBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &TgBot{
		usecases: usecases,
		botTg:    bot,
		ctx:      ctx,
		wg:       wg,
	}, nil
}
