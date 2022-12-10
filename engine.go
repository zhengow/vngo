package vngo

import (
    "github.com/zhengow/vngo/engine"
    "github.com/zhengow/vngo/gateway"
)

func NewBacktestingEngine() *engine.BacktestingEngine {
    return engine.NewBacktestingEngine()
}

func NewLiveTradeEngine(gI gateway.GatewayInterface) *engine.LiveTradeEngine {
    return engine.NewLiveTradeEngine(gI)
}
