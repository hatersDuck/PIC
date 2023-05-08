package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hatersDuck/PIC/pkg/database"
)

func (b *Bot) cbOnTrade(del *tgbotapi.DeleteMessageConfig, cal *tgbotapi.CallbackQuery) (tgbotapi.Chattable, bool) {
	userRow := &database.User{}
	row := b.db.QueryRow("SELECT strategy_id FROM users WHERE user_id = $1", cal.From.ID)
	err := row.Scan(&userRow.StrategyId)

	update_msg := tgbotapi.NewEditMessageText(del.ChatID, del.MessageID, "Торговля успешно запущена")

	if err != nil {
		log.Fatal(err)
	}

	if userRow.StrategyId != 1 {
		_, err := b.db.Exec("UPDATE users SET status_trade = 'Y' WHERE user_id = $1", del.ChatID)
		if err != nil {
			log.Println("Ошибка с включением торговли", err)

			alert := tgbotapi.NewCallback(cal.ID, "Ошибка с включением торговли")
			b.bot.AnswerCallbackQuery(alert)

			return nil, false
		}

		buttons := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnBack, callbackAccount),
			tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnOffTrade, callbackOffTrade),
		))
		update_msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		}

		return update_msg, true
	}

	alert := tgbotapi.NewCallback(cal.ID, b.messages.NoStrategy)
	b.bot.AnswerCallbackQuery(alert)
	return nil, false

}

func (b *Bot) cbOffTradedel(del *tgbotapi.DeleteMessageConfig, cal *tgbotapi.CallbackQuery) (tgbotapi.Chattable, bool) {
	update_msg := tgbotapi.NewEditMessageText(del.ChatID, del.MessageID, "Торговля успешно выключена")

	_, err := b.db.Exec("UPDATE users SET status_trade = 'N' WHERE user_id = $1", cal.From.ID)
	if err != nil {
		log.Println("Ошибка с выключением торговли", err)
		alert := tgbotapi.NewCallback(cal.ID, b.messages.NoStrategy)
		b.bot.AnswerCallbackQuery(alert)
		return nil, false
	}

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnBack, callbackAccount),
		tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnOnTrade, callbackOnTrade),
	))
	update_msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}

	return update_msg, true
}

const countRows = 5

func (b *Bot) cbChangeStrategy(del *tgbotapi.DeleteMessageConfig, strategyId int) (tgbotapi.Chattable, bool) {
	update_msg := tgbotapi.NewEditMessageText(del.ChatID, del.MessageID, b.messages.ChangeStrategy)

	if strategyId != 1 {
		b.db.Exec("UPDATE users SET strategy_id = $1 WHERE user_id = $2", strategyId, del.ChatID)
	}

	userRow := &database.User{}
	row := b.db.QueryRow("SELECT strategy_id FROM users WHERE user_id = $1", del.ChatID)
	row.Scan(&userRow.StrategyId)

	rows, err := b.db.Query("SELECT strategy_id, title FROM tradestrategy WHERE status='active' LIMIT $1", countRows)
	if err != nil {
		log.Println("Ошибка запроса стратегий")
		return nil, false
	}
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, countRows+2)
	for rows.Next() {
		strategyRow := &database.TradeStrategy{}
		if err := rows.Scan(&strategyRow.Id, &strategyRow.Title); err != nil {
			log.Fatal(err)
		}
		title := ""
		if userRow.StrategyId == strategyRow.Id {
			title = "✅"
		}

		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(title+strategyRow.Title, fmt.Sprintf("%s&%d", callbackChangeStrategy, strategyRow.Id)),
		))
	}

	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnBack, callbackAccount),
	))

	update_msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}

	return update_msg, true
}

func (b *Bot) cbApiKeyEmpty(del *tgbotapi.DeleteMessageConfig) tgbotapi.Chattable {
	update_msg := tgbotapi.NewEditMessageText(del.ChatID, del.MessageID, b.messages.AddApi)
	b.db.Exec("UPDATE users SET state_in_bot = 'ap' WHERE user_id = $1", del.ChatID)
	return update_msg
}

func (b *Bot) cbSecretKeyEmpty(del *tgbotapi.DeleteMessageConfig) tgbotapi.Chattable {
	update_msg := tgbotapi.NewEditMessageText(del.ChatID, del.MessageID, b.messages.AddApi)
	b.db.Exec("UPDATE users SET state_in_bot = 'se' WHERE user_id = $1", del.ChatID)
	return update_msg
}

func (b *Bot) cbApiKeyReady(del *tgbotapi.DeleteMessageConfig) tgbotapi.Chattable {
	update_msg := tgbotapi.NewEditMessageText(del.ChatID, del.MessageID, b.messages.DeleteKey)
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnBack, callbackAccount),
		tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnAccept, callbackDeleteApi),
	))
	update_msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
	return update_msg
}

func (b *Bot) cbSecretKeyReady(del *tgbotapi.DeleteMessageConfig) tgbotapi.Chattable {
	update_msg := tgbotapi.NewEditMessageText(del.ChatID, del.MessageID, b.messages.DeleteKey)
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnBack, callbackAccount),
		tgbotapi.NewInlineKeyboardButtonData(b.messages.BtnAccept, callbackDeleteSecret),
	))
	update_msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
	return update_msg
}

func (b *Bot) cbReport(del *tgbotapi.DeleteMessageConfig) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(del.ChatID, "")
	return msg
}

func (b *Bot) cbDeleteApi(del *tgbotapi.DeleteMessageConfig) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(del.ChatID, "")
	return msg
}

func (b *Bot) cbDeleteSecret(del *tgbotapi.DeleteMessageConfig) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(del.ChatID, "")
	return msg
}
