package usecases

import (
	"quizlet_bot/internal/domain/manager/repo_manager"
	"quizlet_bot/internal/domain/models/db_models"
)

type WordsUsecases struct {
	repo repo_manager.WordsRepository
}

func NewWordsUsecases(repo repo_manager.WordsRepository) *WordsUsecases {
	return &WordsUsecases{repo: repo}
}

func (u *WordsUsecases) AddWord(data db_models.Words) error {
	return u.repo.AddWord(data)
}

func (u *WordsUsecases) GetWordsBySet(setName string) ([]string, error) {
	return u.repo.GetWordsBySet(setName)
}

func (u *WordsUsecases) GetTranslationBySet(setName string) ([]string, error) {
	return u.repo.GetTranslationBySet(setName)
}

func (u *WordsUsecases) GetWordsByUser(tgId int64) ([]string, error) {
	return u.repo.GetWordsByUser(tgId)
}

func (u *WordsUsecases) GetTranslationByUser(tgId int64) ([]string, error) {
	return u.repo.GetTranslationByUser(tgId)
}
