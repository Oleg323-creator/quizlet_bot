package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type TgBot struct {
	Ucase *usecases.Usecases
	BotTg *tgbotapi.BotAPI
}

func NewTgBot(ucase *usecases.Usecases, cli *crypto_comp.CryptoCompareAPI, token string) *TgBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logrus.Panic(err)
	}

	return &TgBot{
		Ucase: ucase,
		Cli:   cli,
		BotTg: bot,
	}
}
