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

	selectQuery := sq.Select("id").
		From("users").
		Where(squirrel.Eq{"id": tgId})

	query, args, err := sq.Insert("stats").
		Columns("user_id", "stats").
		Values(selectQuery, 1).
		Suffix("ON CONFLICT (stats) DO UPDATE SET stats = stats.stats + 1").
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

	return nil
}

func (r *StatsPostgres) GetStats(tgId int64) (int64, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Select("stats").
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

	var stats int64
	err = r.db.Pool.QueryRow(context.Background(), query, args...).Scan(&stats)
	if err != nil {
		return 0, err
	}

	return stats, nil
}
