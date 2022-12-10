package vngo

type TemplateStrategy struct {
    BaseStrategy
}

func (s *TemplateStrategy) OnBars(bars map[string]Bar) {
    println("implement me")
}

func (s *TemplateStrategy) UpdateTrade(trade TradeData) {
    println("implement me")
}
