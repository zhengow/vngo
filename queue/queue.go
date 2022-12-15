package queue

import "github.com/zhengow/vngo/models"

type Queue struct {
    Bars     chan map[models.Symbol]models.Bar
    Continue chan struct{}
    Trade    chan models.TradeData
    Order    chan *models.Order
}

func NewQueue() *Queue {
    return &Queue{
        Bars:     make(chan map[models.Symbol]models.Bar),
        Continue: make(chan struct{}),
        Trade:    make(chan models.TradeData),
        Order:    make(chan *models.Order),
    }
}
