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
	err := u.repo.AddTopic(topic, words)
	if err != nil {
		return err
	}
	return nil
}

func (u *TopicsAndWordsUsecases) ChooseTopic(data db_models.Topics) ([]string, error) {
	words, err := u.repo.ChooseTopic(data)
	if err != nil {
		return nil, err
	}

	return words, nil
}
