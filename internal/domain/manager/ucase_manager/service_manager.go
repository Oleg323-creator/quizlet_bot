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

type TopicsUsecases interface {
	AddSet(topic db_models.Sets) error
	SetsList(tgId int64) ([]string, error)
}

type WordsUsecases interface {
	AddWord(data db_models.Words) error
	GetWordsBySet(setName string) ([]string, error)
	GetTranslationBySet(setName string) ([]string, error)
	GetWordsByUser(tgId int64) ([]string, error)
	GetTranslationByUser(tgId int64) ([]string, error)
}

type StatsUsecases interface {
	AddStats(tgId int64) error
	GetStats(tgId int64) (int64, error)
}

type ManagerUsecases struct {
	UsersUsecases
	TopicsUsecases
	WordsUsecases
	StatsUsecases
}

func NewManagerUsecases(repo *repo_manager.ManagerRepo) *ManagerUsecases {
	return &ManagerUsecases{
		UsersUsecases:  usecases.NewUsersUsecases(repo.UsersRepository),
		TopicsUsecases: usecases.NewTopicsUsecases(repo.TopicsRepository),
		WordsUsecases:  usecases.NewWordsUsecases(repo.WordsRepository),
		StatsUsecases:  usecases.NewStatsUsecases(repo.StatsRepository),
	}
}
