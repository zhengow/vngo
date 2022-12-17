package engine

import (
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/queue"
    "time"
)

type BaseEngine struct {
    Kind     string
    symbols  []models.Symbol
    interval models.Interval
    account  models.Account
    queue.Queue
}

//func (b *BaseEngine) addSymbol(name string, exchange models.Exchange) *BaseEngine {
//    symbol := models.NewSymbol(name, exchange)
//    b.symbols = append(b.symbols, symbol)
//    return b
//}

func (b *BaseEngine) AddSymbols(symbols []models.Symbol) *BaseEngine {
    b.symbols = symbols
    return b
}

func (b *BaseEngine) GetSymbols() []models.Symbol {
    return b.symbols
}

func (b *BaseEngine) SetInterval(interval models.Interval) *BaseEngine {
    b.interval = interval
    return b
}

func (b *BaseEngine) GetInterval() models.Interval {
    return b.interval
}

func (b *BaseEngine) LoadHistoryData(_, _ time.Time) {
}

func (b *BaseEngine) Register(q queue.Queue) {
    b.Queue = q
}

func (b *BaseEngine) GetAccount() models.Account {
    return b.account
}

type Engine interface {
    Register(queue.Queue)
    GetAccount() models.Account
}
