package queue

import "github.com/zhengow/vngo/models"

type Queue struct {
    Bars  chan map[string]models.Bar
    wait  chan struct{}
    Trade chan *models.TradeData
    Order chan *models.Order
}

func NewQueue() *Queue {
    return &Queue{
        Bars:  make(chan map[string]models.Bar),
        wait:  make(chan struct{}),
        Trade: make(chan *models.TradeData),
        Order: make(chan *models.Order),
    }
}

func (q *Queue) SendBarsSync(bars map[string]models.Bar) {
    q.Bars <- bars
    <-q.wait
}

func (q *Queue) SendBars(bars map[string]models.Bar) {
    q.Bars <- bars
}

func (q *Queue) SendTradeSync(trade *models.TradeData) {
    q.Trade <- trade
    <-q.wait
}

func (q *Queue) SendTrade(trade *models.TradeData) {
    q.Trade <- trade
}

func (q *Queue) SendOrderSync(order *models.Order) {
    q.Order <- order
    <-q.wait
}

func (q *Queue) SendOrder(order *models.Order) {
    q.Order <- order
}

func (q *Queue) Continue() {
    q.wait <- struct{}{}
}
