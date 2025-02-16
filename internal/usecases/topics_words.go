package usecases

import (
	"quizlet_bot/internal/domain/manager/repo_manager"
	"quizlet_bot/internal/domain/models/db_models"
)

type TopicsAndWordsUsecases struct {
	repo repo_manager.TopicsRepository
}

func NewTopicsAndWordsUsecases(repo repo_manager.TopicsRepository) *TopicsAndWordsUsecases {
	return &TopicsAndWordsUsecases{repo: repo}
}

func (u *TopicsAndWordsUsecases) AddTopic(topic db_models.Sets, words []db_models.Words) error {
	return u.repo.AddTopic(topic, words)
}

func (u *TopicsAndWordsUsecases) WordsBySetName(data db_models.Sets) ([]string, error) {
	return u.repo.WordsBySetName(data)
}

func (u *TopicsAndWordsUsecases) SetsList(tgId int64) ([]string, error) {
	return u.repo.SetsList(tgId)
}
