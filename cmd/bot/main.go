package main

import (
	pocket "github.com/IskanderSh/go-pocket-sdk"
	"github.com/IskanderSh/th_bot_golang/pkg/config"
	"github.com/IskanderSh/th_bot_golang/pkg/repository"
	"github.com/IskanderSh/th_bot_golang/pkg/repository/bolt_db"
	"github.com/IskanderSh/th_bot_golang/pkg/server"
	"github.com/IskanderSh/th_bot_golang/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg)
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := bolt_db.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.AuthServerURL, cfg.Messages)

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, cfg.TelegramBotURL)

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return db, nil
}
