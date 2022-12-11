package gateway

import (
    "context"
    "fmt"
    "github.com/zhengow/vngo/strategy"
    "github.com/zhengow/vngo/types"
    "github.com/zhengow/vngo/utils"
    "log"
    "time"

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

func (f *binanceFutureClient) LoadBarData(symbol strategy.Symbol, interval types.Interval) ([]strategy.Bar, error) {
    res, err := f.client.NewKlinesService().Symbol(symbol.Name).Interval(string(interval)).Do(context.Background())
    if err != nil {
        log.Fatalf("load bar data err: %v", err)
    }
    bars := make([]strategy.Bar, len(res))
    for i, v := range res {
        openValue, highValue, lowValue, closeValue, volumeValue, err := utils.ParseBarData(v.Open, v.High, v.Low, v.Close, v.Volume)
        if err != nil {
            return nil, err
        }
        bars[i] = strategy.Bar{
            Symbol:       symbol,
            Datetime:     strategy.NewVnTime(time.UnixMilli(v.OpenTime)),
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

func (f *binanceFutureClient) WebSocketKLine(symbols []strategy.Symbol, interval types.Interval) {
    wsKlineHandler := func(event *futures.WsKlineEvent) {
        fmt.Println(event)
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
