package vngo

import (
	"github.com/zhengow/vngo/engine"
	"github.com/zhengow/vngo/models"
	"github.com/zhengow/vngo/queue"
	"github.com/zhengow/vngo/strategy"
)

type Loader struct {
	*queue.Queue
	strategy.Strategy
	listenSync bool
}

func (l *Loader) Register(q *queue.Queue, s strategy.Strategy) {
	l.Queue = q
	l.Strategy = s
	l.Strategy.SetQueue(q)
}

func (l *Loader) ListenInit() {
	for bars := range l.Queue.Init.GetChan() {
		l.OnInit(bars)
		l.Queue.Init.Continue()
	}
}

func (l *Loader) ListenBar() {
	for bars := range l.Queue.Bars.GetChan() {
		l.OnBars(bars)
		if l.listenSync {
			l.Queue.Bars.Continue()
		}
	}
}

func (l *Loader) ListenTrade() {
	for trade := range l.Queue.Trade.GetChan() {
		l.UpdateTrade(*trade)
		if l.listenSync {
			l.Queue.Trade.Continue()
		}
	}
}

func (l *Loader) ListenOrder() {
	for order := range l.Queue.Order.GetChan() {
		l.UpdateOrder(order)
		if l.listenSync {
			l.Queue.Order.Continue()
		}
	}
}

func Register(e engine.Engine, s strategy.Strategy) {
	q := queue.NewQueue()
	e.SetQueue(q)
	s.Inject(e.GetAccount())
	l := &Loader{
		listenSync: e.GetKind() == models.EngineEnum.BACKTESTENGINE,
	}
	l.Register(q, s)
    go l.ListenInit()
	go l.ListenBar()
	go l.ListenTrade()
	go l.ListenOrder()
}
