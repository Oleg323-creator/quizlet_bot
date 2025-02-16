package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"quizlet_bot/internal/db"
	"quizlet_bot/internal/domain/models/db_models"
)

type TopicsPostgres struct {
	db *db.WrapperDB
}

func NewTopicsPostgres(db *db.WrapperDB) *TopicsPostgres {
	return &TopicsPostgres{db: db}
}

func (r *TopicsPostgres) AddSet(topic db_models.Sets) error {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Insert("topics").
		Columns("topic", "user_id").
		Select(
			sq.Select("u.id").
				From("users u").
				Where(squirrel.Eq{"u.tg_id": topic.TgId}),
		).
		Suffix("ON CONFLICT (topic, user_id) DO NOTHING").
		ToSql()

	args = append([]interface{}{topic.SetName}, args...)

	if err != nil {
		logrus.Errorf("ERR building query: %v", err)
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
		logrus.Errorf("ERR inserting into stats: %v", err)
		return err
	}

	return nil
}

func (r *TopicsPostgres) SetsList(tgId int64) ([]string, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Select("topic").
		From("topics").
		Join("users ON topics.user_id = users.id").
		Where(squirrel.And{
			squirrel.Eq{"users.tg_id": tgId},
		}).
		ToSql()

	if err != nil {
		logrus.Errorf("ERR creating quiery Get SetName")
		return nil, err
	}

	conn, err := r.db.Pool.Acquire(r.db.Ctx)

	if err != nil {
		logrus.Errorf("ERR acquiring connection: %s", err.Error())
		return nil, err
	}

	defer conn.Release() // Closing conn after req

	rows, execErr := r.db.Pool.Query(r.db.Ctx, query, args...)
	if execErr != nil {
		logrus.Errorf("ERR to execute SQL query: %v", execErr)
		return nil, execErr
	}
	defer rows.Close()

	var topics []string

	for rows.Next() {
		var topic string
		if err = rows.Scan(&topic); err != nil {
			logrus.Errorf("ERR to scan result: %v", err)
			return nil, err
		}

		topics = append(topics, topic)
	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)
	}

	return topics, nil
}
