package usecases

import "quizlet_bot/internal/domain/manager/repo_manager"

type StatsUsecases struct {
	repo repo_manager.StatsRepository
}

func NewStatsUsecases(repo repo_manager.StatsRepository) *StatsUsecases {
	return &StatsUsecases{repo: repo}
}

func (u *StatsUsecases) AddStats(tgId int64) error {
	return u.repo.AddStats(tgId)
}

func (u *StatsUsecases) GetStats(tgId int64) (int64, error) {
	return u.repo.GetStats(tgId)
}
