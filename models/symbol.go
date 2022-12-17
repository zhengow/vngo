package models

import (
    "fmt"
)

type Symbol struct {
    Name     string `gorm:"column:symbol"`
    Exchange Exchange
}

func NewSymbol(name string, exchange Exchange) Symbol {
    return Symbol{
        Name:     name,
        Exchange: exchange,
    }
}

func NewSymbols(names []string, exchanges []Exchange) []Symbol {
    symbols := make([]Symbol, len(names))
    if len(exchanges) == 0 {
        return nil
    }
    for idx, name := range names {
        _exchange := exchanges[0]
        if idx < len(exchanges) {
            _exchange = exchanges[idx]
        }
        symbols[idx] = NewSymbol(name, _exchange)
    }
    return symbols
}

func (s *Symbol) FullName() string {
    return fmt.Sprintf("%s.%s", s.Name, s.Exchange)
}
