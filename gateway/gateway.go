package gateway

import (
    "github.com/zhengow/vngo/strategy"
    "github.com/zhengow/vngo/types"
)

type GatewayInterface interface {
    LoadBarData(strategy.Symbol, types.Interval) ([]strategy.Bar, error)
    WebSocketKLine([]strategy.Symbol, types.Interval)
}
