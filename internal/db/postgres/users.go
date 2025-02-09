package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"quizlet_bot/internal/db"
	"quizlet_bot/internal/domain/models/db_models"
)

type UsersPostgres struct {
	db *db.WrapperDB
}

func NewUsersPostgres(db *db.WrapperDB) *UsersPostgres {
	return &UsersPostgres{db: db}
}

func (r *UsersPostgres) AddUser(data db_models.Users) error {
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

	ex, err := conn.Exec(r.db.Ctx, query, args...)
	if err != nil {
		logrus.Errorf("Error executing query: %s", err.Error())
		return err
	}

	if ex.RowsAffected() > 0 {
		logrus.Infof("NEW USER ADDED! Tg ID: %d", data.TgId)
	}

	_ = r.db.Pool.Stat()

	return nil
}
