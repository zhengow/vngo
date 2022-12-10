package gateway

import "github.com/zhengow/vngo/model"

type GatewayInterface interface {
	LoadBarData(*model.Symbol) ([]model.Bar, error)
}
