package engine

import (
    "github.com/zhengow/vngo/strategy"
    "github.com/zhengow/vngo/utils"
)

type BaseAccount struct {
    Cash      float64
    Orders    map[string]*strategy.Order
    Positions map[strategy.Symbol]float64
    baseOrderRuler
}

func (b *BaseAccount) Buy(symbol strategy.Symbol, price, volume float64) string {
    return ""
}
func (b *BaseAccount) Sell(symbol strategy.Symbol, price, volume float64) string {
    return ""
}
func (b *BaseAccount) CancelAll() {
}

func (b *BaseAccount) CancelById(orderId string) {
}

func (b *BaseAccount) GetPositions() map[strategy.Symbol]float64 {
    return nil
}

func (b *BaseAccount) GetCash() float64 {
    return 0
}

func (b *BaseAccount) GetBalance() float64 {
    return 0
}

func (b *BaseAccount) SetFilters(priceFilter, volumeFilter map[strategy.Symbol]numberFilter) {
    b.PriceFilter = priceFilter
    b.VolumeFilter = volumeFilter
}

type numberFilter struct {
    tickSize  float64
    precision int
}

type baseOrderRuler struct {
    PriceFilter  map[strategy.Symbol]numberFilter
    VolumeFilter map[strategy.Symbol]numberFilter
}

func (b *baseOrderRuler) PriceToTickSize(symbol strategy.Symbol, price float64) float64 {
    if b.PriceFilter == nil {
        return price
    }
    if filter, ok := b.PriceFilter[symbol]; ok {
        return utils.AmountToTickSize(filter.tickSize, filter.precision, price)
    } else {
        return price
    }
}

func (b *baseOrderRuler) VolumeToTickSize(symbol strategy.Symbol, quantity float64) float64 {
    if b.VolumeFilter == nil {
        return quantity
    }
    if filter, ok := b.VolumeFilter[symbol]; ok {
        return utils.AmountToTickSize(filter.tickSize, filter.precision, quantity)
    } else {
        return quantity
    }
}
