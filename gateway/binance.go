package gateway

import (
    "context"
    "fmt"
    "github.com/zhengow/vngo/models"
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

func (f *binanceFutureClient) LoadBarData(symbol *models.Symbol) ([]models.Bar, error) {
    res, err := f.client.NewKlinesService().Symbol(symbol.Name).Interval(string(symbol.Interval)).Do(context.Background())
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
            Symbol:       *symbol,
            Datetime:     *models.NewVnTime(time.UnixMilli(v.OpenTime)),
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

func (f *binanceFutureClient) WebSocketKLine(symbols []*models.Symbol) {
    wsKlineHandler := func(event *futures.WsKlineEvent) {
        fmt.Println(event)
    }
    errHandler := func(err error) {
        fmt.Println(err)
    }
    combined := make(map[string]string)
    for _, symbol := range symbols {
        combined[symbol.Name] = string(symbol.Interval)
    }
    doneC, _, err := futures.WsCombinedKlineServe(combined, wsKlineHandler, errHandler)
    if err != nil {
        fmt.Println(err)
        return
    }
    <-doneC
}
