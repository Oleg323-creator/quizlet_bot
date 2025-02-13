package tg

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
	"github.com/sirupsen/logrus"
	"quizlet_bot/internal/domain/models/db_models"
	"strings"
)

type User struct {
	FSM *fsm.FSM
}

// Хранилище пользователей (map[UserID] -> *User)
var users = make(map[int64]*User)

func (t *TgBot) Bot() (*tgbotapi.BotAPI, error) {

	//go t.Alert(os.Getenv("ADDR"), t.BotTg)

	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Start the bot"},
		{Command: "choose_set", Description: "Choose set"},
		{Command: "create_set", Description: "Create set"},
	}

	config := tgbotapi.NewSetMyCommands(commands...)
	_, err := t.botTg.Request(config)
	if err != nil {
		logrus.Errorf("ERR setting commands: %v", err)
		return nil, err
	}

	//t.botTg.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.botTg.GetUpdatesChan(u)

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

				for update := range updates {
					logrus.Info("UPDATE")
					var chatID int64
					var text string

					if update.Message != nil {
						chatID = update.Message.Chat.ID
						text = update.Message.Text
					} else if update.CallbackQuery != nil {
						chatID = update.CallbackQuery.Message.Chat.ID
						text = update.CallbackQuery.Data
					} else {
						continue
					}

					// Если пользователь новый — создаем для него FSM
					if _, exists := users[chatID]; !exists {
						users[chatID] = t.NewUserFSM()
					}
					user := users[chatID]

					logrus.Info("UPDATE 1")
					// Обрабатываем команды
					switch text {
					case "/start":
						logrus.Info("!!!!!!!")
						t.handleStart(chatID, user)

					default:
						logrus.Info("UPDATE 3")
						t.handleUserResponse(chatID, user, text)
					}
				}
			} else if update.CallbackQuery != nil {
				logrus.Info("UPDATE 4")
				t.handleUserResponse(update.Message.Chat.ID, users[update.Message.Chat.ID], update.CallbackQuery.Data)
			}
		}
	}
}

func (t *TgBot) NewUserFSM() *User {
	return &User{
		FSM: fsm.NewFSM(
			"start", // Начальное состояние
			fsm.Events{
				{Name: "choose_starting_option", Src: []string{"start"}, Dst: "waiting_for_starting_option"},

				{Name: "choose_set", Src: []string{"waiting_for_starting_option"}, Dst: "waiting_for_choosing_set"},
				{Name: "working_with_set", Src: []string{"waiting_for_choosing_set"}, Dst: "working_with_set"},

				{Name: "create_set", Src: []string{"waiting_for_starting_option"}, Dst: "waiting_for_starting_creating"},
				{Name: "enter_set_name", Src: []string{"waiting_for_starting_creating"}, Dst: "waiting_for_entering_name"},
				{Name: "add_word", Src: []string{"waiting_for_entering_name"}, Dst: "waiting_for_adding_word"},

				{Name: "update_set", Src: []string{"waiting_for_starting_option"}, Dst: "waiting_for_set_updating"},

				{Name: "delete_set", Src: []string{"waiting_for_starting_option"}, Dst: "waiting_for_deleting_set"},

				{Name: "complete", Src: []string{"waiting_for_choosing_set", "waiting_for_adding_word",
					"waiting_for_set_updating", "waiting_for_deleting_set"}, Dst: "start"},
			},
			fsm.Callbacks{},
		),
	}
}

func (t *TgBot) handleStart(chatID int64, user *User) {
	user.FSM.Event(t.ctx, "choose_starting_option")

	msg := tgbotapi.NewMessage(chatID, "Choose option")

	btn1 := tgbotapi.NewInlineKeyboardButtonData("Choose set", "choose_set")
	btn2 := tgbotapi.NewInlineKeyboardButtonData("Create set", "create_set")
	btn3 := tgbotapi.NewInlineKeyboardButtonData("Update set", "update_set")
	btn4 := tgbotapi.NewInlineKeyboardButtonData("Delete set", "delete_set")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(btn1, btn2),
		tgbotapi.NewInlineKeyboardRow(btn3, btn4),
	)

	msg.ReplyMarkup = keyboard
	t.botTg.Send(msg)
}

func (t *TgBot) handleUserResponse(chatID int64, user *User, text string) {
	switch {
	case strings.HasPrefix(text, "choose_set"):
		err := t.ChooseSetUserResponse(chatID, user)
		if err != nil {
			logrus.Error(err)
		}

	case strings.HasPrefix(text, "set_was_chosen"):
		err := t.WordsBySetName(chatID, user, text)
		if err != nil {
			logrus.Error(err)
		}

	case strings.HasPrefix(text, "i"):
		err := t.WordsBySetName(chatID, user, text)
		if err != nil {
			logrus.Error(err)
		}

	case strings.HasPrefix(text, "create_set"):
		logrus.Info("create_set")

		user.FSM.Event(t.ctx, "waiting_for_starting_creating")

		_, err := t.usecases.SetsList(chatID)
		if err != nil {
			logrus.Errorf("ERR choosing topic")
			msg := tgbotapi.NewMessage(chatID, "There is no such set ")
			t.botTg.Send(msg)
		}

		msg := tgbotapi.NewMessage(chatID, "")
		t.botTg.Send(msg)

	case strings.HasPrefix(text, "update_set"):
		logrus.Info("update_set")
		user.FSM.Event(t.ctx, "ask_age") // Переход в ожидание возраста
		msg := tgbotapi.NewMessage(chatID, "Сколько тебе лет?")
		t.botTg.Send(msg)

	case strings.HasPrefix(text, "delete_set"):
		logrus.Info("delete_set")
		user.FSM.Event(t.ctx, "ask_age") // Переход в ожидание возраста
		msg := tgbotapi.NewMessage(chatID, "Сколько тебе лет?")
		t.botTg.Send(msg)

	default:
		logrus.Info("default")
		msg := tgbotapi.NewMessage(chatID, "Я не понял тебя. Напиши /start.")
		t.botTg.Send(msg)
	}
}

func (t *TgBot) ChooseSetUserResponse(chatID int64, user *User) error {
	logrus.Info("choose_set")
	user.FSM.Event(t.ctx, "waiting_for_choosing_set")

	topics, err := t.usecases.SetsList(chatID)
	if err != nil {
		logrus.Errorf("ERR choosing topic")
		msg := tgbotapi.NewMessage(chatID, "There is no such set ")
		t.botTg.Send(msg)
	}

	var keyboardSlice [][]tgbotapi.InlineKeyboardButton
	var btnsInRowSlice []tgbotapi.InlineKeyboardButton

	if len(topics) == 0 {
		logrus.Errorf("You don't have any sets")
		return err
	}
	for _, topic := range topics {
		btnsInRowSlice = append(btnsInRowSlice, tgbotapi.NewInlineKeyboardButtonData(topic, fmt.Sprintf("set_was_chosen"+topic)))

		if len(btnsInRowSlice) > 3 {
			keyboardSlice = append(keyboardSlice, btnsInRowSlice)
			btnsInRowSlice = nil // Очищаем строку кнопок
		}
	}

	if len(btnsInRowSlice) > 0 {
		keyboardSlice = append(keyboardSlice, btnsInRowSlice)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardSlice...)

	msg := tgbotapi.NewMessage(chatID, "Choose set:")

	msg.ReplyMarkup = keyboard
	t.botTg.Send(msg)

	return nil
}

func (t *TgBot) WordsBySetName(chatID int64, user *User, callback string) error {
	logrus.Info("working_with_set")

	var setName string
	/*
		if strings.HasPrefix(callback, "i_know") {
			setName = strings.TrimPrefix(callback, "i_know")
		} else if strings.HasPrefix(callback, "i_don't_know") {
			setName = strings.TrimPrefix(callback, "i_don't_know")
		}
	*/

	setName = strings.TrimPrefix(callback, "set_was_chosen")

	user.FSM.Event(t.ctx, "working_with_set")

	data := db_models.Topics{Topic: setName, TgId: chatID}
	words, err := t.usecases.WordsBySetName(data)
	if err != nil {
		return err
	}

	for _, word := range words {
		btn1 := tgbotapi.NewInlineKeyboardButtonData("I know", fmt.Sprintf("i_know"+word))
		btn2 := tgbotapi.NewInlineKeyboardButtonData("I don't know", fmt.Sprintf("i_don't_know"+word))

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(btn1, btn2),
		)

		msg := tgbotapi.NewMessage(chatID, word)

		msg.ReplyMarkup = keyboard
		t.botTg.Send(msg)
	}

	return err
}

/*				case "choose_set":
					logrus.Info("choose_set")
					user.FSM.Event(t.ctx, "waiting_for_choosing_set")

					topics, err := t.usecases.SetsList(chatID)
					if err != nil {
						logrus.Errorf("ERR choosing topic")
						msg := tgbotapi.NewMessage(chatID, "There is no such set ")
						t.botTg.Send(msg)
					}

					var keyboardSlice [][]tgbotapi.InlineKeyboardButton
					var btnsInRowSlice []tgbotapi.InlineKeyboardButton

					counter := 0
					if len(topics) == 0 {
						logrus.Errorf("LENTH IS 0")
						continue
					}
					for _, topic := range topics {
						btnsInRowSlice = append(btnsInRowSlice, tgbotapi.NewInlineKeyboardButtonData(topic, topic))

						if len(btnsInRowSlice) > 3 {
							keyboardSlice = append(keyboardSlice, btnsInRowSlice)
							btnsInRowSlice = nil // Очищаем строку кнопок
							counter++
						}
					}

					keyboardSlice = append(keyboardSlice, btnsInRowSlice)

					keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardSlice...)

					msg := tgbotapi.NewMessage(chatID, "Choose set:")

					msg.ReplyMarkup = keyboard
					t.botTg.Send(msg)

				case "help":
					user.FSM.Event(t.ctx, "help")
					t.botTg.Send(tgbotapi.NewMessage(chatID, "Вот инструкция по использованию бота..."))
*/

/*
	f := fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: "choose_set", Src: []string{"start"}, Dst: "waiting_for_choosing_set"},

			{Name: "create_set", Src: []string{"start"}, Dst: "waiting_for_starting_creating"},
			{Name: "enter_set_name", Src: []string{"waiting_for_starting_creating"}, Dst: "waiting_for_entering_name"},
			{Name: "add_word", Src: []string{"waiting_for_entering_name"}, Dst: "waiting_for_adding_word"},
			{Name: "finish_creating", Src: []string{"waiting_for_adding_word"}, Dst: "start"},

			{Name: "update_set", Src: []string{"start"}, Dst: "waiting_for_set_updating"},

			{Name: "delete_set", Src: []string{"start"}, Dst: "waiting_for_deleting_set"},
		},
		fsm.Callbacks{},
	)

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

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Available commands:\n"+
							"/choose_set - Choose set\n/create_set - Create set")
						t.botTg.Send(msg)
					case "choose_set":
						logrus.Info("Got /choose_set command")

						userStates[update.Message.Chat.ID] = f

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose set")
						t.botTg.Send(msg)
					case "create_set":
						logrus.Info("Got /create_set command")

						userStates[update.Message.Chat.ID] = "creating_set" //saiving status

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Enter set name")
						t.botTg.Send(msg)

					default:
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
						t.botTg.Send(msg)
					}
				} else if userState, exists := userStates[update.Message.Chat.ID]; exists {
					if userState == "choosing_set" {
						set := update.Message.Text
						logrus.Infof("Set chousen: %s", set)

						data := db_models.Topics{
							Topic: set,
							TgId:  update.Message.From.ID,
						}

						_, err := t.usecases.WordsBySetName(data)
						if err != nil {
							logrus.Info("")
							return nil, err
						}

						/*	if len(addr) != 34 {
							confirmationMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Incorrect address")
							t.botTg.Send(confirmationMsg)
							continue
						}*/
/*
	} else if userState == "creating_set" {
		//set := update.Message.Text

		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Add words to the set(word-translate)")

		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Add word"),
				tgbotapi.NewKeyboardButton("Finish creating"),
				tgbotapi.NewKeyboardButton("Cancel"),
			))
		msg.ReplyMarkup = keyboard
		t.botTg.Send(msg)

		userStates[update.Message.Chat.ID] = "adding words"
		if userState == "adding words" {

		}

	}
*/
/*	stats, err := t.StatsForTg(addr)
							if err != nil {
								logrus.Infof("ERR getting stats: fro tg")
							}

							confirmationMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your stats:\n "+stats)
							t.BotTg.Send(confirmationMsg)
	/*
							delete(userStates, update.Message.Chat.ID)*/
