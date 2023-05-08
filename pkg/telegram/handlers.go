package telegram

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hatersDuck/PIC/pkg/database"
)

func (b *Bot) handlerUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		switch {
		case update.Message != nil:
			if update.Message.IsCommand() {
				b.handleCommand(update.Message)
			} else {
				b.handleMessage(update.Message)
			}

		case update.CallbackQuery != nil:
			go b.handleCallback(update.CallbackQuery)

		default:
			log.Printf("Undefiend update %d", update.UpdateID)
		}
	}
}

const (
	callbackOnTrade        = "on_trade"
	callbackOffTrade       = "off_trade"
	callbackChangeStrategy = "change_strategy"
	callbackApiKeyEmpty    = "api_key_empty"
	callbackSecretKeyEmpty = "secret_key_empty"
	callbackApiKeyReady    = "api_key_ready"
	callbackSecretKeyReady = "secret_key_ready"
	callbackReport         = "report"

	callbackDeleteApi    = "delete_api_key"
	callbackDeleteSecret = "delete_secret_key"
	callbackAccount      = "account"
)

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) error {
	deleteConfig := tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	send := true
	var msg tgbotapi.Chattable
	data := strings.Split(callback.Data, "&")

	switch data[0] {
	case callbackOnTrade:
		msg, send = b.cbOnTrade(&deleteConfig, callback)

	case callbackOffTrade:
		msg, send = b.cbOffTradedel(&deleteConfig, callback)
	case callbackChangeStrategy:
		if len(data) > 1 {
			if strategyID, err := strconv.Atoi(data[1]); err == nil {
				msg, send = b.cbChangeStrategy(&deleteConfig, strategyID)
			} else {
				log.Println("WTF")
			}
		} else {
			msg, send = b.cbChangeStrategy(&deleteConfig, 1)
		}

	case callbackApiKeyEmpty:
		msg = b.cbApiKeyEmpty(&deleteConfig)

	case callbackSecretKeyEmpty:
		msg = b.cbSecretKeyEmpty(&deleteConfig)

	case callbackApiKeyReady:
		msg = b.cbApiKeyReady(&deleteConfig)

	case callbackSecretKeyReady:
		msg = b.cbSecretKeyReady(&deleteConfig)

	case callbackDeleteApi:
		b.db.Exec("UPDATE users SET api_key = 'empty' WHERE user_id = $1", callback.From.ID)

		b.bot.Send(tgbotapi.NewMessage(int64(callback.From.ID), "Ключ успешно удалён"))
		b.bot.Send(deleteConfig)
		b.accountCommand(int64(callback.From.ID), *callback.From)

		send = false

	case callbackDeleteSecret:
		b.db.Exec("UPDATE users SET secret_key = 'empty' WHERE user_id = $1", callback.From.ID)

		b.bot.Send(tgbotapi.NewMessage(int64(callback.From.ID), "Ключ успешно удалён"))
		b.bot.Send(deleteConfig)
		b.accountCommand(int64(callback.From.ID), *callback.From)

		send = false

	case callbackReport:
		cmd := exec.Command("python3", "scripts/create_diagram.py", fmt.Sprintf("%d", callback.From.ID))
		stdout, err := cmd.Output()

		filename := string(stdout)

		if err != nil {
			fmt.Println(err.Error())
			b.bot.Send(tgbotapi.NewMessage(int64(callback.From.ID), err.Error()))
		} else {
			photoBytes, err := ioutil.ReadFile(filename)
			if err != nil {
				panic(err)
			}
			photoFileBytes := tgbotapi.FileBytes{
				Name:  "diagram",
				Bytes: photoBytes,
			}

			message := tgbotapi.NewPhotoUpload(int64(callback.From.ID), photoFileBytes)

			send = false
			_, err = b.bot.Send(message)
			fmt.Println(err)
		}

	case callbackAccount:
		b.accountCommand(int64(callback.From.ID), *callback.From)
		send = false
		b.bot.Send(deleteConfig)
	}
	if send {
		b.bot.Send(msg)
	}
	return nil
}

func (b *Bot) handleState(del *tgbotapi.DeleteMessageConfig) tgbotapi.Chattable {
	return nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	deleteConfig := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	userRow := &database.User{}
	row := b.db.QueryRow("SELECT state_in_bot FROM users WHERE user_id = $1", message.From.ID)
	err := row.Scan(&userRow.StateInBot)
	if err != nil {
		log.Println("WTF failed state")
	}
	apiKeyRegex := regexp.MustCompile("^[A-Z0-9a-z]{64}$")
	var msg tgbotapi.MessageConfig

	switch userRow.StateInBot {
	//todo Сделать рефакторинг
	case "ap":
		if apiKeyRegex.MatchString(message.Text) {
			//todo Надо добавить шифрофку
			b.db.Exec("UPDATE users SET api_key = $1 WHERE user_id = $2", message.Text, message.From.ID)
			msg = tgbotapi.NewMessage(message.Chat.ID, "api key успешно добавлен. Сообщение удалено в целях безопасности")

			b.bot.Send(msg)
			b.bot.Send(deleteConfig)

			b.accountCommand(message.Chat.ID, *message.From)
		} else {
			msg = tgbotapi.NewMessage(message.Chat.ID, "Это не api_key, попробуйте ещё раз. Для выхода нажмите /start")
			b.bot.Send(msg)
		}
	case "se":
		if apiKeyRegex.MatchString(message.Text) {
			//todo Надо добавить шифрофку
			b.db.Exec("UPDATE users SET secret_key = $1 WHERE user_id = $2", message.Text, message.From.ID)
			msg = tgbotapi.NewMessage(message.Chat.ID, "secret key успешно добавлен. Сообщение удалено в целях безопасности")

			b.bot.Send(msg)
			b.bot.Send(deleteConfig)

			b.accountCommand(message.Chat.ID, *message.From)
		} else {
			msg = tgbotapi.NewMessage(message.Chat.ID, "Это не secret key, попробуйте ещё раз. Для выхода нажмите /start")
			b.bot.Send(msg)
		}
	default:
		msg = tgbotapi.NewMessage(message.Chat.ID, "Бот в основном управляется кнопками под сообщениями пожалуйста не пишите. Если у вас нет меню нажмите /start")
		b.bot.Send(msg)
	}

	return nil
}

const (
	commandStart   = "start"
	commandAccount = "account"
)

func (b *Bot) handleCommand(command *tgbotapi.Message) error {
	sw := command.Command()

	switch sw {
	case commandStart:
		b.startCommand(command.Chat.ID, *command.From)
	case commandAccount:
		b.accountCommand(command.Chat.ID, *command.From)
	default:
		b.accountCommand(command.Chat.ID, *command.From)
	}
	return nil
}
