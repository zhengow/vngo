package strategy

import (
    "github.com/zhengow/vngo/model"
)

type TemplateStrategy struct {
    BaseStrategy
}

func (s *TemplateStrategy) OnBars(bars map[string]model.Bar) {
    println("implement me")
}

func (s *TemplateStrategy) UpdateTrade(trade model.TradeData) {
    println("implement me")
}
