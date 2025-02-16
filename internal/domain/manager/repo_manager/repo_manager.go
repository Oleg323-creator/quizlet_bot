package repo_manager

import (
	"quizlet_bot/internal/db"
	"quizlet_bot/internal/db/postgres"
	"quizlet_bot/internal/domain/models/db_models"
)

type Migrator interface {
	Up() error
	Down() error
}

type UsersRepository interface {
	AddUser(data db_models.Users) error
}

type TopicsRepository interface {
	AddSet(topic db_models.Sets) error
	SetsList(tgId int64) ([]string, error)
}

type WordsRepository interface {
	AddWord(data db_models.Words) error
	GetWordsBySet(setName string) ([]string, error)
	GetTranslationBySet(setName string) ([]string, error)
	GetWordsByUser(tgId int64) ([]string, error)
	GetTranslationByUser(tgId int64) ([]string, error)
}

type StatsRepository interface {
	AddStats(tgId int64) error
	GetStats(tgId int64) (int64, error)
}

type ManagerRepo struct {
	Migrator
	UsersRepository
	TopicsRepository
	WordsRepository
	StatsRepository
}

func NewManagerRepo(db *db.WrapperDB) *ManagerRepo {
	return &ManagerRepo{
		Migrator:         postgres.NewMigratorPostgres(db),
		UsersRepository:  postgres.NewUsersPostgres(db),
		TopicsRepository: postgres.NewTopicsPostgres(db),
		WordsRepository:  postgres.NewWordsPostgres(db),
		StatsRepository:  postgres.NewStatsPostgres(db),
	}
}
