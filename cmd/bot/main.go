package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hatersDuck/PIC/config"
	"github.com/hatersDuck/PIC/pkg/telegram"
	"github.com/hatersDuck/PIC/pkg/trade"
	"github.com/jackc/pgx"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	conn, err := pgx.Connect(pgx.ConnConfig{
		Database: "pic",
		Password: cfg.DatabasePassword,
		User:     "danila",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	connTrade, err := pgx.Connect(pgx.ConnConfig{
		Database: "pic",
		Password: cfg.DbPasswordTrade,
		User:     "trade_binance",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close()
	defer connTrade.Close()

	errChan := make(chan error)
	trade := *trade.NewTrade(cfg.TimeSleep, cfg.TestNet, connTrade)
	go trade.Start(errChan)

	go func() {
		if err := telegram.NewBot(bot, cfg.Messages, conn).Start(); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		select {
		case err := <-errChan:
			log.Printf("Error: %v", err)
		}
	}

}
