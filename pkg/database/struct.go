package database

import "time"

type (
	TradeStrategy struct {
		Id     int    `db:"strategy_id"`
		Title  string `db:"title"`
		Path   string `db:"path_file"`
		Status string `db:"status"`
		Symbol string `db:"symbol"`
	}

	HistoryData struct {
		Id     int       `db:"history_id"`
		Date   time.Time `db:"date"`
		Symbol string    `db:"symbol"`
		Price  float32   `db:"price"`
	}

	User struct {
		Id         int64  `db:"user_id"`
		ApiKey     string `db:"api_key"`
		SecretKey  string `db:"secret_key"`
		StrategyId int    `db:"strategy_id"`
		Status     rune   `db:"status_trade"`
		Username   string `db:"username"`
		StateInBot string `db:"state_in_bot"`
		Succes     bool   `db:"succes"`
	}

	Orders struct {
		Id         int       `db:"order_id"`
		Type       string    `db:"type"`
		Date       time.Time `db:"date"`
		StrategyId int       `db:"strategy_id"`
		HistoryId  int       `db:"history_id"`
	}

	Transaction struct {
		Id       int       `db:"transaction_id"`
		OrderId  int       `db:"order_id"`
		UserID   int64     `db:"user_id"`
		Date     time.Time `db:"date"`
		Quantity float32   `db:"quantity"`
	}
)
