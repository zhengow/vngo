package vngo

import (
    "github.com/zhengow/vngo/engine"
    "github.com/zhengow/vngo/queue"
    "github.com/zhengow/vngo/strategy"
)

type Loader struct {
    queue.Queue
    strategy.Strategy
}

func (l *Loader) Register(q queue.Queue, s strategy.Strategy) {
    l.Queue = q
    l.Strategy = s
}

func (l *Loader) ListenBar() {
    for bars := range l.Queue.Bars {
        l.OnBars(bars)
        l.Continue()
    }
}

func (l *Loader) ListenTrade() {
    for trade := range l.Queue.Trade {
        l.UpdateTrade(*trade)
        l.Continue()
    }
}

func (l *Loader) ListenOrder() {
    for order := range l.Queue.Order {
        l.UpdateOrder(order)
        l.Continue()
    }
}

func Register(e engine.Engine, s strategy.Strategy) {
    q := queue.NewQueue()
    e.Register(*q)
    s.Inject(e.GetAccount())
    l := &Loader{}
    l.Register(*q, s)
    go l.ListenBar()
    go l.ListenTrade()
    go l.ListenOrder()
}
