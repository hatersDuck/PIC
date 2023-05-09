package trade

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/hatersDuck/PIC/pkg/database"
	"github.com/jackc/pgx"
)

type Trade struct {
	timeSleep int
	testNet   bool

	db *pgx.Conn
}

func NewTrade(timeSleep int, testNet bool, conn *pgx.Conn) *Trade {
	return &Trade{
		timeSleep: timeSleep,
		testNet:   testNet,
		db:        conn,
	}
}

func (t *Trade) Start(errChan chan<- error) {
	binance.UseTestnet = t.testNet

	log.Println("[Start Trade]")
	for {
		rows, err := t.db.Query("SELECT strategy_id, path_file, symbol FROM tradeStrategy WHERE status = 'active'")
		log.Println("get strategy")
		if err != nil {
			log.Println(err)
			errChan <- err
		}

		rows_strategy := make([]database.TradeStrategy, 0, 5)
		for rows.Next() {
			row := &database.TradeStrategy{}
			rows.Scan(&row.Id, &row.Path, &row.Symbol)
			rows_strategy = append(rows_strategy, *row)
		}

		for _, row := range rows_strategy {
			row.Symbol = strings.ReplaceAll(row.Symbol, " ", "")
			log.Println("Start", row.Path)
			cmd := exec.Command("python3", row.Path)
			stdout, _ := cmd.Output()
			log.Println("Trigger =", string(stdout))
			res, _ := binance.NewClient("", "").NewListPricesService().Symbol(row.Symbol).Do(context.Background())
			t.db.Exec("INSERT INTO historyData(symbol, price) VALUES($1, $2)", row.Symbol, res[0].Price)
			data := &database.HistoryData{}
			t.db.QueryRow("SELECT history_id FROM historyData ORDER BY history_id DESC LIMIT 1").Scan(&data.Id)
			t.db.Exec("INSERT INTO orders(type, history_id, strategy_id) VALUES($1, $2, $3)", string(stdout), data.Id, row.Id)

			data_order := &database.Orders{}
			t.db.QueryRow("SELECT order_id FROM orders ORDER BY order_id DESC LIMIT 1").Scan(&data_order.Id)

			//todo тут должна быть дешифровка
			users, _ := t.db.Query("select user_id, api_key, secret_key from users where status_trade = 'Y' and strategy_id = $1", row.Id)

			rows_users := make([]database.User, 0, 128)
			for users.Next() {
				user := &database.User{}
				users.Scan(&user.Id, &user.ApiKey, &user.SecretKey)
				rows_users = append(rows_users, *user)
			}

			for _, user := range rows_users {
				client := &ClientBinance{
					cl: *binance.NewClient(user.ApiKey, user.SecretKey),
					id: user.Id,
				}
				switch string(stdout) {
				case "BUY":
					ball, _ := client.Balance(row.Symbol[3:])

					bl, _ := strconv.ParseFloat(ball, 64)
					st, _ := strconv.ParseFloat(res[0].Price, 64)
					quantity := fmt.Sprintf("%f", bl*0.5/st)

					_, err := client.CreateOrder(row.Symbol, res[0].Price, binance.SideTypeBuy, quantity)
					if err == nil {
						t.db.Exec("INSERT INTO transactions(order_id, user_id) VALUES($1, $2, $3)", data_order.Id, client.id)
					}

				case "SELL":
					ball, _ := client.Balance(row.Symbol[:3])
					bl, _ := strconv.ParseFloat(ball, 64)
					quantity := fmt.Sprintf("%f", bl*0.5)

					client.CreateOrder(row.Symbol, res[0].Price, binance.SideTypeSell, quantity)
					if err == nil {
						t.db.Exec("INSERT INTO transactions(order_id, user_id, quantity) VALUES($1, $2, $3)", data_order.Id, client.id, quantity)
					} else {
						errChan <- err
					}
				}
			}
		}
		time.Sleep(time.Duration(t.timeSleep) * time.Second)
	}
}
