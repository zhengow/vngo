package vngo

type Symbol struct {
    Symbol   string
    Exchange Exchange
    Interval Interval
}

func NewSymbol(symbol string, exchange Exchange, interval Interval) *Symbol {
    return &Symbol{
        Symbol:   symbol,
        Exchange: exchange,
        Interval: interval,
    }
}
