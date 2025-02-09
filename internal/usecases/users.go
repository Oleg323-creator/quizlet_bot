package usecases

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"quizlet_bot/internal/domain/manager/repo_manager"
	"quizlet_bot/internal/domain/models/db_models"
)

type UsersUsecases struct {
	repo repo_manager.UsersRepository
}

func NewUsersUsecases(repo repo_manager.UsersRepository) *UsersUsecases {
	return &UsersUsecases{repo: repo}
}

func (u *UsersUsecases) AddUser(user *tgbotapi.User) error {
	dbUser := db_models.Users{
		TgId:      user.ID,
		Username:  user.UserName,
		Firstname: user.FirstName,
		LastName:  user.LastName,
	}

	err := u.repo.AddUser(dbUser)
	if err != nil {
		return err
	}

	return nil
}
