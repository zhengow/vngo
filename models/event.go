package models

type OrderEvent struct {
    Direction string
    Symbol
    Price  float64
    Volume float64
}

type CancelEvent struct{}
