package gateway

import (
    "github.com/adshao/go-binance/v2"
    "github.com/adshao/go-binance/v2/futures"
)

func NewClient() *futures.Client {
    futuresClient := binance.NewFuturesClient("", "")
    return futuresClient
}
