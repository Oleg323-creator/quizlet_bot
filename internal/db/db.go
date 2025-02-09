package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type ConnectionConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

type WrapperDB struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
}

func NewDB(ctx context.Context, cfg ConnectionConfig) (*WrapperDB, error) {
	var pool *pgxpool.Pool
	var err error
	pool, err = NewPostgresDB(ctx, cfg)
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return nil, err
	}

	return &WrapperDB{
		Pool: pool,
		Ctx:  ctx,
	}, nil
}

func (db *WrapperDB) Close() {
	db.Pool.Close()
}
