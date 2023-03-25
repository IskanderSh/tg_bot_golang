package telegram

import (
	pocket "github.com/IskanderSh/go-pocket-sdk"
	"github.com/IskanderSh/th_bot_golang/pkg/config"
	"github.com/IskanderSh/th_bot_golang/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string

	messages config.Messages
}

func NewBot(bot *tgbotapi.BotAPI, pocketClient *pocket.Client, tr repository.TokenRepository, redirectURL string, messages config.Messages) *Bot {
	return &Bot{bot: bot, pocketClient: pocketClient, tokenRepository: tr, redirectURL: redirectURL, messages: messages}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.HandleError(update.Message.Chat.ID, err)
			}
			continue
		}
		if err := b.handleMessage(update.Message); err != nil {
			b.HandleError(update.Message.Chat.ID, err)
		}
	}
}
