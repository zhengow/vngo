package gateway

import (
	"github.com/zhengow/vngo/models"
	"github.com/zhengow/vngo/queue"
)

type GatewayInterface interface {
	LoadBarData(models.Symbol, models.Interval) ([]models.Bar, error)
	WebSocketKLine([]models.Symbol, models.Interval, *queue.Queue)
}
