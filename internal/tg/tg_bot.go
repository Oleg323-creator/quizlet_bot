package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"quizlet_bot/internal/domain/models/db_models"
)

func (t *TgBot) Bot() (*tgbotapi.BotAPI, error) {

	//go t.Alert(os.Getenv("ADDR"), t.BotTg)

	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Start the bot"},
		{Command: "choose_topic", Description: "Choose topic"},
		{Command: "create_topic", Description: "Create topic"},
	}

	config := tgbotapi.NewSetMyCommands(commands...)
	_, err := t.botTg.Request(config)
	if err != nil {
		logrus.Errorf("ERR setting commands: %v", err)
		return nil, err
	}

	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60
	updates := t.botTg.GetUpdatesChan(upd)
	userStates := make(map[int64]string) //for saiving status of bot

	for {
		select {
		case <-t.ctx.Done():
			logrus.Info("Bot stopped due to context cancellation")
			return nil, err
		case update := <-updates:
			if update.Message != nil {

				err := t.usecases.AddUser(update.Message.From)
				if err != nil {
					logrus.Errorf("ERR adding uses: %v", err)
					return nil, err
				}

				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "start":
						logrus.Info("Got /start command")

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Available commands:\n/start - "+
							"Start the bot\n/choose_topic - Choose topic\n/create_topic - Create topic")
						t.botTg.Send(msg)
					case "choose_topic":
						logrus.Info("Got /choose_topic command")

						userStates[update.Message.Chat.ID] = "choosing_topic"

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose topic")
						t.botTg.Send(msg)
					case "create_topic":
						logrus.Info("Got /create_topic command")

						userStates[update.Message.Chat.ID] = "creating_topic" //saiving status

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Enter topic name")
						t.botTg.Send(msg)

					default:
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
						t.botTg.Send(msg)
					}
				} else if userState, exists := userStates[update.Message.Chat.ID]; exists {
					if userState == "choosing_topic" {
						topic := update.Message.Text
						logrus.Infof("Topic chousen: %s", topic)

						data := db_models.Topics{
							Topic: topic,
							TgId:  update.Message.From.ID,
						}

						_, err := t.usecases.ChooseTopic(data)
						if err != nil {
							logrus.Info("")
							return nil, err
						}

						/*	if len(addr) != 34 {
							confirmationMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Incorrect address")
							t.botTg.Send(confirmationMsg)
							continue
						}*/

					} else if userState == "creating_topic" {
						//topic := update.Message.Text

						msg := tgbotapi.NewMessage(update.Message.Chat.ID,
							"Add words to the topic(word-translate)")

						keyboard := tgbotapi.NewReplyKeyboard(
							tgbotapi.NewKeyboardButtonRow(
								tgbotapi.NewKeyboardButton("Add word"),
								tgbotapi.NewKeyboardButton("Create topic"),
								tgbotapi.NewKeyboardButton("Cancel"),
							))
						msg.ReplyMarkup = keyboard
						t.botTg.Send(msg)

						userStates[update.Message.Chat.ID] = "adding words"
						if userState == "adding words" {

						}

					}

					/*	stats, err := t.StatsForTg(addr)
						if err != nil {
							logrus.Infof("ERR getting stats: fro tg")
						}

						confirmationMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your stats:\n "+stats)
						t.BotTg.Send(confirmationMsg)

						delete(userStates, update.Message.Chat.ID)*/
				}
			}
		}
	}
}
