package postgres

import (
	"context"
	"github.com/jackc/tern/migrate"
	"github.com/sirupsen/logrus"
	"quizlet_bot/internal/db"
)

var MigrationsDirectory = "./migrations"

type MigratorPostgres struct {
	db *db.WrapperDB
}

func NewMigratorPostgres(db *db.WrapperDB) *MigratorPostgres {
	return &MigratorPostgres{db: db}
}

func (r *MigratorPostgres) Up() error {
	conn, err := r.db.Pool.Acquire(r.db.Ctx)
	defer conn.Release()
	if err != nil {
		logrus.WithError(err).Error("Failed to acquire DB connection")
		return err
	}
	migrator, err := migrate.NewMigrator(context.Background(), conn.Conn(), "schema_version")
	if err != nil {
		logrus.WithError(err).Error("Failed to create migrator")
		return err
	}

	logrus.Info("Loading migrations...")
	if err = migrator.LoadMigrations(MigrationsDirectory); err != nil {
		logrus.WithError(err).Error("Failed to load migrations")
		return err
	}

	logrus.Info("Running migrations...")
	if err = migrator.Migrate(context.Background()); err != nil {
		logrus.WithError(err).Error("Migration failed")
		return err
	}

	logrus.Info("Migration completed successfully")
	return nil
}

func (r *MigratorPostgres) Down() error {
	return nil
}
