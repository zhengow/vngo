package vngo

type Symbol struct {
    Name     string `gorm:"column:symbol"`
    Exchange Exchange
    Interval Interval
}

func NewSymbol(name string, exchange Exchange, interval Interval) *Symbol {
    return &Symbol{
        Name:     name,
        Exchange: exchange,
        Interval: interval,
    }
}
