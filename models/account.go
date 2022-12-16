package models

import "github.com/zhengow/vngo/utils"

type Account interface {
    Buy(symbol Symbol, price, volume float64) string
    Sell(symbol Symbol, price, volume float64) string
    CancelAll()
    CancelById(orderId string)
    GetPositions() map[Symbol]float64
    GetCash() float64
    GetBalance() float64
}

type BaseAccount struct {
    Cash      float64
    Orders    map[string]*Order
    Positions map[Symbol]float64
    orderFilters
}

func (b *BaseAccount) Buy(symbol Symbol, price, volume float64) string {
    return ""
}
func (b *BaseAccount) Sell(symbol Symbol, price, volume float64) string {
    return ""
}
func (b *BaseAccount) CancelAll() {
}

func (b *BaseAccount) CancelById(orderId string) {
}

func (b *BaseAccount) GetPositions() map[Symbol]float64 {
    return nil
}

func (b *BaseAccount) GetCash() float64 {
    return 0
}

func (b *BaseAccount) GetBalance() float64 {
    return 0
}

func (b *BaseAccount) SetFilters(priceFilter, volumeFilter map[Symbol]numberFilter) {
    b.PriceFilter = priceFilter
    b.VolumeFilter = volumeFilter
}

type numberFilter struct {
    tickSize  float64
    precision int
}

type orderFilters struct {
    PriceFilter  map[Symbol]numberFilter
    VolumeFilter map[Symbol]numberFilter
}

func (b *orderFilters) PriceToTickSize(symbol Symbol, price float64) float64 {
    if b.PriceFilter == nil {
        return price
    }
    if filter, ok := b.PriceFilter[symbol]; ok {
        return utils.AmountToTickSize(filter.tickSize, filter.precision, price)
    } else {
        return price
    }
}

func (b *orderFilters) VolumeToTickSize(symbol Symbol, quantity float64) float64 {
    if b.VolumeFilter == nil {
        return quantity
    }
    if filter, ok := b.VolumeFilter[symbol]; ok {
        return utils.AmountToTickSize(filter.tickSize, filter.precision, quantity)
    } else {
        return quantity
    }
}
