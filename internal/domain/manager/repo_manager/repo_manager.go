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

type TopicsAndWordsRepository interface {
	AddTopic(topic db_models.Topics, words []db_models.Words) error
	WordsBySetName(data db_models.Topics) ([]string, error)
	SetsList(tgId int64) ([]string, error)
}

type StatsRepository interface {
	AddStats(tgId int64) error
	GetStats(tgId int64) (int64, error)
}

type ManagerRepo struct {
	Migrator
	UsersRepository
	TopicsAndWordsRepository
	StatsRepository
}

func NewManagerRepo(db *db.WrapperDB) *ManagerRepo {
	return &ManagerRepo{
		Migrator:                 postgres.NewMigratorPostgres(db),
		UsersRepository:          postgres.NewUsersPostgres(db),
		TopicsAndWordsRepository: postgres.NewTopicsAndWordsPostgres(db),
		StatsRepository:          postgres.NewStatsPostgres(db),
	}
}
