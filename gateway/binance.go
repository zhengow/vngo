package gateway

import (
	"context"
	"fmt"
	"log"

	"github.com/zhengow/vngo/models"
	"github.com/zhengow/vngo/queue"
	"github.com/zhengow/vngo/utils"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
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

func (f *binanceFutureClient) LoadBarData(symbol models.Symbol, interval models.Interval) ([]models.Bar, error) {
	res, err := f.client.NewKlinesService().Symbol(symbol.Name).Interval(string(interval)).Do(context.Background())
	if err != nil {
		log.Fatalf("load bar data err: %v", err)
	}
	bars := make([]models.Bar, len(res))
	for i, v := range res {
		openValue, highValue, lowValue, closeValue, volumeValue, err := utils.ParseBarData(v.Open, v.High, v.Low, v.Close, v.Volume)
		if err != nil {
			return nil, err
		}
		bars[i] = models.Bar{
			Symbol:       symbol,
			Datetime:     models.NewVnTimeByTimestamp(v.OpenTime),
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

func (f *binanceFutureClient) WebSocketKLine(symbols []models.Symbol, interval models.Interval, q *queue.Queue) {
	wsKlineHandler := func(event *futures.WsKlineEvent) {
		// fmt.Println(event)
		if event.Event == "kline" {
			v := event.Kline
			openValue, highValue, lowValue, closeValue, volumeValue, err := utils.ParseBarData(v.Open, v.High, v.Low, v.Close, v.Volume)
			if err != nil {
				log.Fatalf("get websocket kline err: %v", err)
				return
			}
			symbol := models.NewSymbol(event.Symbol, models.ExchangeEnum.BINANCE)
			bar := models.Bar{
				Symbol:       symbol,
				Datetime:     models.NewVnTimeByTimestamp(event.Time),
				Volume:       volumeValue,
				OpenInterest: 0,
				OpenPrice:    openValue,
				HighPrice:    highValue,
				LowPrice:     lowValue,
				ClosePrice:   closeValue,
			}
			q.Bars.Send(map[string]models.Bar{symbol.FullName(): bar})
		}
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	combined := make(map[string]string)
	for _, symbol := range symbols {
		combined[symbol.Name] = string(interval)
	}
	doneC, _, err := futures.WsCombinedKlineServe(combined, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}
