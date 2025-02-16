package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"quizlet_bot/internal/db"
)

type StatsPostgres struct {
	db *db.WrapperDB
}

func NewStatsPostgres(db *db.WrapperDB) *StatsPostgres {
	return &StatsPostgres{db: db}
}

func (r *StatsPostgres) AddStats(tgId int64) error {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	selectQuery, args, err := sq.Select("id").
		From("users").
		Where(squirrel.Eq{"tg_id": tgId}).
		ToSql()

	if err != nil {
		logrus.Errorf("ERR selecting user: %s", err.Error())
		return err
	}

	conn, err := r.db.Pool.Acquire(r.db.Ctx)
	if err != nil {
		logrus.Errorf("Error acquiring connection: %s", err.Error())
		return err
	}
	defer conn.Release()

	var userID int64
	err = conn.QueryRow(r.db.Ctx, selectQuery, args...).Scan(&userID)
	if err != nil {
		logrus.Errorf("User not found or error selecting: %s", err.Error())
		return err
	}

	query, args, err := sq.Insert("stats").
		Columns("user_id", "stat").
		Values(userID, 1).
		Suffix("ON CONFLICT (user_id) DO UPDATE SET stat = stats.stat + 1").
		ToSql()

	if err != nil {
		logrus.Errorf("ERR inserting into DB: %s", err.Error())
		return err
	}

	_, err = conn.Exec(r.db.Ctx, query, args...)
	if err != nil {
		logrus.Errorf("Error executing query: %s", err.Error())
		return err
	}

	return nil
}

func (r *StatsPostgres) GetStats(tgId int64) (int64, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Select("stat", "user_id").
		From("stats").
		Join("users ON stats.user_id = users.id").
		Where(squirrel.And{
			squirrel.Eq{"users.tg_id": tgId},
		}).
		ToSql()

	if err != nil {
		return 0, err
	}

	conn, err := r.db.Pool.Acquire(context.Background())

	if err != nil {
		logrus.Errorf("Error acquiring connection: %s", err.Error())
		return 0, err
	}

	defer conn.Release() // Closing conn after req

	var stat, userId int64
	err = r.db.Pool.QueryRow(context.Background(), query, args...).Scan(&stat, &userId)
	if err != nil {
		return 0, err
	}

	updateQuery, updateArgs, err := sq.Update("stats").
		Set("stat", 0).
		Where(squirrel.Eq{"user_id": userId}).
		ToSql()

	if err != nil {
		logrus.Errorf("Error preparing update query: %s", err.Error())
		return 0, err
	}

	_, err = conn.Exec(context.Background(), updateQuery, updateArgs...)
	if err != nil {
		logrus.Errorf("Error executing update query: %s", err.Error())
		return 0, err
	}

	return stat, nil
}
