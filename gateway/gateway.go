package gateway

import "github.com/zhengow/vngo/models"

type GatewayInterface interface {
    LoadBarData(*models.Symbol) ([]models.Bar, error)
    WebSocketKLine([]*models.Symbol)
}
