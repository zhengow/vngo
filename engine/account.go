package engine

import (
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/utils"
)

type BaseAccount struct {
    Cash      float64
    Orders    map[string]*models.Order
    Positions map[models.Symbol]float64
    baseOrderRuler
}

func (b *BaseAccount) Buy(symbol models.Symbol, price, volume float64) string {
    return ""
}
func (b *BaseAccount) Sell(symbol models.Symbol, price, volume float64) string {
    return ""
}
func (b *BaseAccount) CancelAll() {
}

func (b *BaseAccount) CancelById(orderId string) {
}

func (b *BaseAccount) GetPositions() map[models.Symbol]float64 {
    return nil
}

func (b *BaseAccount) GetCash() float64 {
    return 0
}

func (b *BaseAccount) GetBalance() float64 {
    return 0
}

func (b *BaseAccount) SetFilters(priceFilter, volumeFilter map[models.Symbol]numberFilter) {
    b.PriceFilter = priceFilter
    b.VolumeFilter = volumeFilter
}

type numberFilter struct {
    tickSize  float64
    precision int
}

type baseOrderRuler struct {
    PriceFilter  map[models.Symbol]numberFilter
    VolumeFilter map[models.Symbol]numberFilter
}

func (b *baseOrderRuler) PriceToTickSize(symbol models.Symbol, price float64) float64 {
    if b.PriceFilter == nil {
        return price
    }
    if filter, ok := b.PriceFilter[symbol]; ok {
        return utils.AmountToTickSize(filter.tickSize, filter.precision, price)
    } else {
        return price
    }
}

func (b *baseOrderRuler) VolumeToTickSize(symbol models.Symbol, quantity float64) float64 {
    if b.VolumeFilter == nil {
        return quantity
    }
    if filter, ok := b.VolumeFilter[symbol]; ok {
        return utils.AmountToTickSize(filter.tickSize, filter.precision, quantity)
    } else {
        return quantity
    }
}

type Account struct {
}
