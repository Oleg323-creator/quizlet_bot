package usecases

import (
	"quizlet_bot/internal/domain/manager/repo_manager"
	"quizlet_bot/internal/domain/models/db_models"
)

type TopicsUsecases struct {
	repo repo_manager.TopicsRepository
}

func NewTopicsUsecases(repo repo_manager.TopicsRepository) *TopicsUsecases {
	return &TopicsUsecases{repo: repo}
}

func (u *TopicsUsecases) AddSet(topic db_models.Sets) error {
	return u.repo.AddSet(topic)
}

func (u *TopicsUsecases) SetsList(tgId int64) ([]string, error) {
	return u.repo.SetsList(tgId)
}
