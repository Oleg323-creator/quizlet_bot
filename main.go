package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"quizlet_bot/internal/db"
	"quizlet_bot/internal/domain/manager/repo_manager"
	"quizlet_bot/internal/domain/manager/ucase_manager"
	"quizlet_bot/internal/tg"
	"sync"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file")
	}

	cfg := db.ConnectionConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  "disable",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()

	dbConn, err := db.NewDB(ctx, cfg)
	if err != nil {
		panic(err)
	}
	repo := repo_manager.NewManagerRepo(dbConn)
	usec := ucase_manager.NewManagerUsecases(repo)

	err = repo.Up()
	if err != nil {
		logrus.Info(err)
		return
	}

	wg := &sync.WaitGroup{}

	tgBot, err := tg.NewTgBot(usec, os.Getenv("BOT_TOKEN"), ctx, wg)
	if err != nil {
		logrus.Error(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := tgBot.Bot()
		if err != nil {
			logrus.Errorf("Error running bot: %v", err)
		}

	}()

	wg.Wait()
}
