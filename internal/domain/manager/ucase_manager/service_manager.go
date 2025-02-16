package ucase_manager

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"quizlet_bot/internal/domain/manager/repo_manager"
	"quizlet_bot/internal/domain/models/db_models"
	"quizlet_bot/internal/usecases"
)

type UsersUsecases interface {
	AddUser(user *tgbotapi.User) error
}

type TopicsAndWordsUsecases interface {
	AddTopic(topic db_models.Sets, words []db_models.Words) error
	WordsBySetName(data db_models.Sets) ([]string, error)
	SetsList(tgId int64) ([]string, error)
}

type StatsUsecases interface {
	AddStats(tgId int64) error
	GetStats(tgId int64) (int64, error)
}

type ManagerUsecases struct {
	UsersUsecases
	TopicsAndWordsUsecases
	StatsUsecases
}

func NewManagerUsecases(repo *repo_manager.ManagerRepo) *ManagerUsecases {
	return &ManagerUsecases{
		UsersUsecases:          usecases.NewUsersUsecases(repo.UsersRepository),
		TopicsAndWordsUsecases: usecases.NewTopicsAndWordsUsecases(repo.TopicsRepository),
		StatsUsecases:          usecases.NewStatsUsecases(repo.StatsRepository),
	}
}
