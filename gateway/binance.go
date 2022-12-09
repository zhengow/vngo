package gateway

import (
    "github.com/adshao/go-binance/v2"
    "github.com/adshao/go-binance/v2/futures"
)

func NewFutureClient(apiKey, secretKey string) *futures.Client {
    futuresClient := binance.NewFuturesClient(apiKey, secretKey)
    return futuresClient
}
