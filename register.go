package vngo

import (
    "github.com/zhengow/vngo/engine"
    "github.com/zhengow/vngo/queue"
    "github.com/zhengow/vngo/strategy"
)

type Loader struct {
    *queue.Queue
    strategy.Strategy
}

func (l *Loader) Register(q *queue.Queue, s strategy.Strategy) {
    l.Queue = q
    l.Strategy = s
    l.Strategy.SetQueue(q)
}

func (l *Loader) ListenBar() {
    for bars := range l.Queue.Bars.GetChan() {
        l.OnBars(bars)
        l.Queue.Bars.Continue()
    }
}

func (l *Loader) ListenTrade() {
    for trade := range l.Queue.Trade.GetChan() {
        l.UpdateTrade(*trade)
        l.Queue.Trade.Continue()
    }
}

func (l *Loader) ListenOrder() {
    for order := range l.Queue.Order.GetChan() {
        l.UpdateOrder(order)
        l.Queue.Order.Continue()
    }
}

func Register(e engine.Engine, s strategy.Strategy) {
    q := queue.NewQueue()
    e.SetQueue(q)
    s.Inject(e.GetAccount())
    l := &Loader{}
    l.Register(q, s)
    go l.ListenBar()
    go l.ListenTrade()
    go l.ListenOrder()
}
