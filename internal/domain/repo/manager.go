package repo

import (
	"quizlet_bot/internal/db"
	"quizlet_bot/internal/db/postgres"
	"quizlet_bot/internal/domain/models"
)

type Migrator interface {
	Up() error
	Down() error
}

type Users interface {
	AddUser(data models.Users) error
}

type TopicsAndWords interface {
	AddTopic(topic models.Topics, words []models.Words) error
	GetTopic(data models.Topics) ([]string, error)
}

type ManagerRepo struct {
	Migrator
	Users
	TopicsAndWords
}

func NewManagerRepo(db *db.WrapperDB) *ManagerRepo {
	return &ManagerRepo{
		Migrator:       postgres.NewMigratorPostgres(db),
		Users:          postgres.NewUsersPostgres(db),
		TopicsAndWords: postgres.NewMTopicsPostgres(db),
	}
}
