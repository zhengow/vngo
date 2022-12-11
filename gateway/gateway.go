package gateway

import "github.com/zhengow/vngo/strategy"

type GatewayInterface interface {
    LoadBarData(strategy.Symbol) ([]strategy.Bar, error)
    WebSocketKLine([]strategy.Symbol)
}
