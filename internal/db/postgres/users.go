package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"quizlet_bot/internal/db"
	"quizlet_bot/internal/domain/models"
)

type UserPostgres struct {
	db *db.WrapperDB
}

func NewUsersPostgres(db *db.WrapperDB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) AddUser(data models.Users) error {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Insert("users").
		Columns("tg_id", "username", "first_name", "last_name").
		Values(data.TgId, data.Username, data.Firstname, data.LastName).
		Suffix("ON CONFLICT (tg_id) DO NOTHING").
		ToSql()

	if err != nil {
		logrus.Errorf("ERR insearting into DB: %s", err.Error())
		return err
	}

	conn, err := r.db.Pool.Acquire(r.db.Ctx)

	if err != nil {
		logrus.Errorf("Error acquiring connection: %s", err.Error())
		return err
	}

	defer conn.Release() // Closing conn after req

	_, err = conn.Exec(r.db.Ctx, query, args...)
	if err != nil {
		logrus.Errorf("Error executing query: %s", err.Error())
		return err
	}

	_ = r.db.Pool.Stat()

	logrus.Infof("NEW USER! Tg ID: %s", data.TgId)
	return nil
}
