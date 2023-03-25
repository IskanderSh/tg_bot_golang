package telegram

import (
	"context"
	pocket "github.com/IskanderSh/go-pocket-sdk"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/url"
)

const (
	commandStart = "start"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	_, err := url.ParseQuery(message.Text)
	if err != nil {
		return errInvalidURL
	}
	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return errUnauthorized
	}
	if err = b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		return errUnableToSave
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.SavedSuccessfully)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.AlreadyAuthorized)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.UnknownCommand)
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
