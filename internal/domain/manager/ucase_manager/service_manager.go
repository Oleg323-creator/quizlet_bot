package ucase_manager

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"quizlet_bot/internal/domain/manager/repo_manager"
	"quizlet_bot/internal/domain/models/db_models"
	"quizlet_bot/internal/usecases"
)

type Users interface {
	AddUser(user *tgbotapi.User) error
}

type TopicsAndWords interface {
	AddTopic(topic db_models.Topics, words []db_models.Words) error
	ChooseTopic(data db_models.Topics) ([]string, error)
	TopicsList(tgId int64) ([]string, error)
}

type ManagerUsecases struct {
	Users
	TopicsAndWords
}

func NewManagerUsecases(repo *repo_manager.ManagerRepo) *ManagerUsecases {
	return &ManagerUsecases{
		Users:          usecases.NewUsersUsecases(repo.UsersRepository),
		TopicsAndWords: usecases.NewTopicsAndWordsUsecases(repo.TopicsAndWordsRepository),
	}
}
