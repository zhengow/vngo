package engine

import (
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/queue"
    "time"
)

type BaseEngine struct {
    kind     models.EngineType
    symbols  []models.Symbol
    interval models.Interval
    account  models.Account
    *queue.Queue
}

func NewBaseEngine(kind models.EngineType) *BaseEngine {
    return &BaseEngine{
        kind: kind,
    }
}

func (b *BaseEngine) GetKind() models.EngineType {
    return b.kind
}

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

func (b *BaseEngine) SetQueue(q *queue.Queue) {
    b.Queue = q
}

func (b *BaseEngine) GetAccount() models.Account {
    return b.account
}

type Engine interface {
    SetQueue(*queue.Queue)
    GetAccount() models.Account
}
