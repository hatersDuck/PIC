package trade

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
)

type ClientBinance struct {
	cl binance.Client
	id int64
}

func (c *ClientBinance) Status() bool {
	res, _ := c.cl.NewGetAPIKeyPermission().Do(context.Background())
	return res.EnableSpotAndMarginTrading
}

func (c *ClientBinance) Balance(asset string) (string, error) {
	res, err := c.cl.NewGetAccountService().Do(context.Background())
	if err != nil {
		return "0", err
	}

	for _, balance := range res.Balances {
		if balance.Asset == asset {
			rt := balance.Free
			return rt, nil
		}
	}

	return "0", fmt.Errorf("не удалось найти баланс для актива %s", asset)
}

func (c *ClientBinance) CreateOrder(symbol string, price string, side binance.SideType, quantity string) (*binance.CreateOrderResponse, error) {
	newOrder, err := c.cl.NewCreateOrderService().
		Symbol(symbol).
		Side(side).
		Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(quantity).
		Price(price).Do(context.Background())

	return newOrder, err
}
