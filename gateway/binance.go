package gateway

import (
	"context"
	"fmt"
	"log"
	"time"

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

func (f *binanceFutureClient) LoadBarDataByMinute(symbols []models.Symbol, q *queue.Queue) {
	bars := make(map[string]models.Bar, len(symbols))
	for _, symbol := range symbols {
		bars[symbol.FullName()] = models.Bar{
			Symbol:   symbol,
			Interval: models.IntervalEnum.MINUTE,
		}
	}
	WebSocketKLine(symbols, models.IntervalEnum.MINUTE, bars)
	ticker := time.NewTicker(time.Millisecond * 500)
	for range ticker.C {
		log.Println("tick")
		if time.Now().Second() == 0 {
			if isValidBars(bars) {
				q.Bars.Send(bars)
			}
			ticker.Stop()
			time.Sleep(time.Second * 1)
			ticker.Reset(time.Millisecond * 500)
		}
	}
}

func isValidBars(bars map[string]models.Bar) bool {
	for _, bar := range bars {
		if bar.Datetime.IsZero() {
			return false
		}
	}
	return true
}

func WebSocketKLine(symbols []models.Symbol, interval models.Interval, bars map[string]models.Bar) {
	name2FullName := make(map[string]string, len(symbols))
	for _, symbol := range symbols {
		name2FullName[symbol.Name] = symbol.FullName()
	}
	wsKlineHandler := func(event *futures.WsKlineEvent) {
		if event.Event == "kline" {
			v := event.Kline
			openValue, highValue, lowValue, closeValue, volumeValue, err := utils.ParseBarData(v.Open, v.High, v.Low, v.Close, v.Volume)
			if err != nil {
				log.Fatalf("get websocket kline err: %v", err)
				return
			}
			fullName := name2FullName[event.Symbol]
			bar := bars[fullName]
			bar.SetDatetime(models.NewVnTimeByTimestamp(event.Time))
			bar.SetOpenPrice(openValue).SetHighPrice(highValue).SetLowPrice(lowValue).SetClosePrice(closeValue).SetVolume(volumeValue)
			bars[fullName] = bar
		}
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	combined := make(map[string]string)
	for _, symbol := range symbols {
		combined[symbol.Name] = string(interval)
	}
	_, _, err := futures.WsCombinedKlineServe(combined, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	// <-doneC
}
