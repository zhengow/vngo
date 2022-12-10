package gateway

import (
	"context"
	"log"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/zhengow/vngo/model"
	"github.com/zhengow/vngo/utils"
)

type binanceFutureClient struct {
	client *futures.Client
}

func NewFutureClient(apiKey, secretKey string) *binanceFutureClient {
	futuresClient := binance.NewFuturesClient(apiKey, secretKey)
	return &binanceFutureClient{
		client: futuresClient,
	}
}

func (f *binanceFutureClient) LoadBarData(symbol *model.Symbol) ([]model.Bar, error) {
	res, err := f.client.NewKlinesService().Symbol(symbol.Symbol).Interval(string(symbol.Interval)).Do(context.Background())
	if err != nil {
		log.Fatalf("load bar data err: %v", err)
	}
	bars := make([]model.Bar, len(res))
	for i, v := range res {
		openValue, highValue, lowValue, closeValue, volumeValue, err := utils.ParseBarData(v.Open, v.High, v.Low, v.Close, v.Volume)
		if err != nil {
			return nil, err
		}
		bars[i] = model.Bar{
			Symbol:       *symbol,
			Datetime:     utils.DatabaseTime(time.UnixMilli(v.OpenTime)),
			Interval:     symbol.Interval,
			Volume:       volumeValue,
			OpenInterest: 0,
			OpenPrice:    openValue,
			HighPrice:    highValue,
			LowPrice:     lowValue,
			ClosePrice:   closeValue,
		}
	}
	return bars, nil
}
