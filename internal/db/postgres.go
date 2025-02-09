package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func NewPostgresDB(ctx context.Context, cfg ConnectionConfig) (*pgxpool.Pool, error) {

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host,
		cfg.Port, cfg.DBName, cfg.SSLMode)

	logrus.Infof("%s", connString)
	conf, err := pgxpool.ParseConfig(connString) // Using environment variables instead of a connection string.
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return nil, err
	}

	conf.ConnConfig.LogLevel = pgx.LogLevelWarn
	conf.MaxConns = 50
	conf.ConnConfig.PreferSimpleProtocol = true

	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return nil, err
	}

	if err = getConnection(ctx, pool); err != nil {
		logrus.Errorf("%s", err.Error())
		return nil, err
	}

	return pool, nil
}

// get connection from pool and release
func getConnection(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)

	defer conn.Release()

	if err != nil {
		logrus.Errorf("%s", err.Error())
		return err
	}
	if err = conn.Ping(ctx); err != nil {
		logrus.Errorf("%s", err.Error())
		return err
	}
	return nil
}
