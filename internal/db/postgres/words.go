package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"log"
	"quizlet_bot/internal/db"
	"quizlet_bot/internal/domain/models/db_models"
)

type WordsPostgres struct {
	db *db.WrapperDB
}

func NewWordsPostgres(db *db.WrapperDB) *WordsPostgres {
	return &WordsPostgres{db: db}
}

func (r *WordsPostgres) AddWord(data db_models.Words) error {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Insert("words").
		Columns("word", "translation", "topic_id", "user_id").
		Select(
			squirrel.Select(
				data.Word,
				data.Translate,
				"(SELECT id FROM users WHERE tg_id = ?)",
				"(SELECT id FROM topic WHERE topic = ?)",
			),
		).
		ToSql()

	if err != nil {
		logrus.Errorf("ERR building query for stats: %v", err)
		return err
	}

	args = append(args, data.TgId, data.SetName)

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

func (r *TopicsPostgres) GetWordsBySet(setName string) ([]string, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Select("word").
		From("words").
		Join("topics ON words.topic_id = topics.id").
		Where(squirrel.And{
			squirrel.Eq{"topics.topic": setName},
		}).
		ToSql()

	/*
		query, args, err := sq.Select("word").
		From("words").
		Join("topics ON words.topic_id = topics.id").
		Join("users ON topics.user_id = users.id").
		Where(squirrel.And{
			squirrel.Eq{"topics.topic": data.SetName},
			squirrel.Eq{"users.tg_id": data.TgId},
		}).
		ToSql()
	*/

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
		return nil, err
	}
	defer rows.Close()

	var words []string
	for rows.Next() {
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

func (r *TopicsPostgres) GetTranslationBySet(setName string) ([]string, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Select("translation").
		From("words").
		Join("topics ON words.topic_id = topics.id").
		Where(squirrel.And{
			squirrel.Eq{"topics.topic": setName},
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
		return nil, err
	}
	defer rows.Close()

	var translations []string
	for rows.Next() {
		var translation string
		if err = rows.Scan(&translation); err != nil {
			logrus.Errorf("ERR to scan result: %v", err)
			return nil, err
		}

		translations = append(translations, translation)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return translations, nil
}

func (r *TopicsPostgres) GetWordsByUser(tgId int64) ([]string, error) {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Select("word").
		From("words").
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
		return nil, err
	}
	defer rows.Close()

	var words []string
	for rows.Next() {
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

func (r *TopicsPostgres) GetTranslationByUser(tgId int64) ([]string, error) {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Select("translation").
		From("words").
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
		return nil, err
	}
	defer rows.Close()

	var translations []string
	for rows.Next() {
		var translation string
		if err = rows.Scan(&translation); err != nil {
			logrus.Errorf("ERR to scan result: %v", err)
			return nil, err
		}

		translations = append(translations, translation)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return translations, nil
}
