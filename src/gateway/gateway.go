package gateway

import (
    "github.com/zhengow/vngo"
)

type GatewayInterface interface {
    LoadBarData(*vngo.Symbol) ([]vngo.Bar, error)
    WebSocketKLine([]*vngo.Symbol)
}
