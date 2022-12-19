package queue

import "github.com/zhengow/vngo/models"

type Queue struct {
    Bars  *WrappedChan[map[string]models.Bar]
    Init  *WrappedChan[map[string]models.Bar]
    Trade *WrappedChan[*models.TradeData]
    Order *WrappedChan[*models.Order]
    //test WrappedChan[int]
}

type WrappedChan[T any] struct {
    ch     chan T
    closed bool
    wait   *chan struct{}
}

func (c *WrappedChan[T]) GetChan() chan T {
    return c.ch
}

func (c *WrappedChan[T]) CloseChan() {
    c.closed = true
    close(c.ch)
}

func (c *WrappedChan[T]) Send(data T) {
    if c.closed {
        return
    }
    c.ch <- data
}

func (c *WrappedChan[T]) SendSync(data T) {
    if c.closed {
        return
    }
    c.ch <- data
    <-*c.wait
}

func (c *WrappedChan[T]) SetWaitChan(wait *chan struct{}) {
    c.wait = wait
}

func (c *WrappedChan[T]) Continue() {
    *c.wait <- struct{}{}
}

func NewQueue() *Queue {
    wait := make(chan struct{})
    q := &Queue{
        Bars:  &WrappedChan[map[string]models.Bar]{ch: make(chan map[string]models.Bar)},
        Init:  &WrappedChan[map[string]models.Bar]{ch: make(chan map[string]models.Bar)},
        Trade: &WrappedChan[*models.TradeData]{ch: make(chan *models.TradeData)},
        Order: &WrappedChan[*models.Order]{ch: make(chan *models.Order)},
    }
    q.Bars.SetWaitChan(&wait)
    q.Init.SetWaitChan(&wait)
    q.Trade.SetWaitChan(&wait)
    q.Order.SetWaitChan(&wait)
    return q
}
