package gateway

import (
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/types"
)

type GatewayInterface interface {
    LoadBarData(models.Symbol, types.Interval) ([]models.Bar, error)
    WebSocketKLine([]models.Symbol, types.Interval)
}
