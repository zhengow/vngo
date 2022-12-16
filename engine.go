package vngo

import (
    "github.com/zhengow/vngo/engine"
    "github.com/zhengow/vngo/engine/backtesting"
    "github.com/zhengow/vngo/gateway"
)

func NewBacktestingEngine() *backtesting.BacktestingEngine {
    return backtesting.NewBacktestingEngine()
}

func NewLiveTradeEngine(gI gateway.GatewayInterface) *engine.LiveTradeEngine {
    return engine.NewLiveTradeEngine(gI)
}
