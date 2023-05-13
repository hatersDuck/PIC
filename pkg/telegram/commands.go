package telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hatersDuck/PIC/pkg/database"
)

func (b *Bot) startCommand(chat_id int64, user tgbotapi.User) {
	message := tgbotapi.NewMessage(chat_id, b.messages.Start)
	b.bot.Send(message)
	b.accountCommand(chat_id, user)
}

func (b *Bot) accountCommand(chat_id int64, user tgbotapi.User) {
	userRow := &database.User{}
	row := b.db.QueryRow("SELECT * FROM users WHERE user_id = $1", user.ID)

	err := row.Scan(&userRow.Id, &userRow.ApiKey, &userRow.SecretKey, &userRow.StrategyId, &userRow.Status, &userRow.Username, &userRow.StateInBot, &userRow.Succes)

	if err != nil {
		log.Println("New user ", user.ID, err)
		if userRow.Id != int64(user.ID) {
			_, err := b.db.Exec("insert into users (user_id, username, strategy_id) values ($1, $2, 1)", user.ID, user.UserName)
			if err != nil {
				log.Println("EXEC failed ", chat_id, user.ID, err)
			}
		}
		userRow.Status = 'N'
		userRow.ApiKey = "empty" + strings.Repeat(" ", 59)
		userRow.SecretKey = "empty" + strings.Repeat(" ", 59)
	}
	if userRow.StateInBot != "no" {
		b.db.Exec("UPDATE users SET state_in_bot = 'no' WHERE user_id = $1", user.ID)
	}

	message := b.createAccount(chat_id, userRow)
	b.bot.Send(message)
}

func (b *Bot) createAccount(chat_id int64, user *database.User) tgbotapi.MessageConfig {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, 3)

	emptyStr := "empty" + strings.Repeat(" ", 59)
	msgTex := b.messages.Account
	msgWar := b.messages.Warnings
	status := b.messages.OffTrade
	no_api_keys := b.messages.NoApiKeys

	if user.Status == 'N' {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnOnTrade, callbackOnTrade),
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnChangeStrategy, "change_strategy"),
		))
	} else {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnOffTrade, callbackOffTrade),
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnChangeStrategy, "change_strategy"),
		))
	}

	if user.ApiKey == emptyStr && user.SecretKey == emptyStr {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnApiKeyEmpty, "api_key_empty"),
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnSecretKeyEmpty, "secret_key_empty"),
		))
	} else if user.ApiKey != emptyStr && user.SecretKey == emptyStr {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnApiKeyReady, "api_key_ready"),
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnSecretKeyEmpty, "secret_key_empty"),
		))
	} else if user.ApiKey == emptyStr && user.SecretKey != emptyStr {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnApiKeyEmpty, "api_key_empty"),
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnSecretKeyReady, "secret_key_ready"),
		))
	} else {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnApiKeyReady, "api_key_ready"),
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnSecretKeyReady, "secret_key_ready"),
		))
		if user.Status == 'Y' && user.Succes {
			status = b.messages.OnTrade
		}

		if !user.Succes {
			msgWar = fmt.Sprintf(msgWar, b.messages.NoSuccess)
		}
		no_api_keys = ""

	}
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnReport, "report"),
	))

	if msgWar == b.messages.Warnings && no_api_keys != "" {
		msgWar = fmt.Sprintf(msgWar, no_api_keys)
	} else if msgWar == b.messages.Warnings {
		msgWar = ""
	}
	log.Println(msgWar)
	message := tgbotapi.NewMessage(chat_id, fmt.Sprintf(msgTex, status)+msgWar)
	message.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: buttons}
	message.ParseMode = "Markdown"
	return message
}
