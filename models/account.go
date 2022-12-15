package models

type Account interface {
    Buy(symbol Symbol, price, volume float64) string
    Sell(symbol Symbol, price, volume float64) string
    CancelAll()
    CancelById(orderId string)
    GetPositions() map[Symbol]float64
    GetCash() float64
    GetBalance() float64
}
