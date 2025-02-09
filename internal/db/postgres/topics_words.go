package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"log"
	"quizlet_bot/internal/db"
	"quizlet_bot/internal/domain/models"
)

type TopicsPostgres struct {
	db *db.WrapperDB
}

func NewMTopicsPostgres(db *db.WrapperDB) *TopicsPostgres {
	return &TopicsPostgres{db: db}
}

func (r *TopicsPostgres) AddTopic(topic models.Topics, words []models.Words) error {

	//using TX to add data to both tables(topics and words)

	tx, err := r.db.Pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		logrus.Errorf("ERR starting transaction: %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(context.Background())
		} else if err != nil {
			tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Insert("topics").
		Columns("topic", "user_id").
		Select(
			sq.Select(sq.Expr("?", topic.Topic), "u.id").
				From("users u").
				Where(squirrel.Eq{"u.tg_id": topic.TgId}).
				Suffix("ON CONFLICT (topic, user_id) DO NOTHING"),
		).
		ToSql()

	if err != nil {
		logrus.Errorf("ERR building query: %v", err)
		return err
	}

	var topicId uint64
	err = tx.QueryRow(r.db.Ctx, query, args...).Scan(&topicId)
	if err != nil {
		logrus.Errorf("ERR inserting address: %v", err)
		return err
	}

	for _, word := range words {
		query, args, err = sq.Insert("words").
			Columns("word", "translate", "topic_id").
			Values(word.Word, word.Translate, topicId).
			Suffix("ON CONFLICT (word, translate, topic_id) DO NOTHING").
			ToSql()

		if err != nil {
			logrus.Errorf("ERR building query for stats: %v", err)
			return err
		}
	}

	_, err = tx.Exec(r.db.Ctx, query, args...)
	if err != nil {
		logrus.Errorf("ERR inserting into stats: %v", err)
		return err
	}

	return nil
}

func (r *TopicsPostgres) GetTopic(data models.Topics) ([]string, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Select("*").
		From("words").
		Join("topics ON words.topic_id = topics.id").
		Where(squirrel.And{
			squirrel.Eq{"topics.topic": data.Topic},
			squirrel.Eq{"topics.tg_id": data.TgId},
		}).
		ToSql()

	if err != nil {
		logrus.Errorf("ERR creating quiery Get Topic")
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
		return nil, err
	}
	defer rows.Close()

	var words []string
	if rows.Next() {
		var word string
		if err = rows.Scan(&word); err != nil {
			logrus.Errorf("ERR to scan result: %v", err)
			return nil, err
		}

		words = append(words, word)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return words, nil
}

/*
func (r *TopicsPostgres) AddTopic(data models.Topics) error {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Insert("topics").
		Columns("topic", "user_id").
		Select(
			sq.Select(sq.Expr("?", data.Topic), "u.id").
				From("users u").
				Where(squirrel.Eq{"u.tg_id": data.TgId}),
		).
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

	logrus.Infof("NEW USER! Tg ID: %d", data.TgId)
	return nil
}
*/
