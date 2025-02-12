package usecases

import (
	"quizlet_bot/internal/domain/manager/repo_manager"
	"quizlet_bot/internal/domain/models/db_models"
)

type TopicsAndWordsUsecases struct {
	repo repo_manager.TopicsAndWordsRepository
}

func NewTopicsAndWordsUsecases(repo repo_manager.TopicsAndWordsRepository) *TopicsAndWordsUsecases {
	return &TopicsAndWordsUsecases{repo: repo}
}

func (u *TopicsAndWordsUsecases) AddTopic(topic db_models.Topics, words []db_models.Words) error {
	return u.repo.AddTopic(topic, words)
}

func (u *TopicsAndWordsUsecases) ChooseTopic(data db_models.Topics) ([]string, error) {
	return u.repo.ChooseTopic(data)
}

func (u *TopicsAndWordsUsecases) TopicsList(tgId int64) ([]string, error) {
	return u.repo.TopicsList(tgId)
}
