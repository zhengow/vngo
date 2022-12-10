package vngo

import (
    "github.com/zhengow/vngo/engine"
    "github.com/zhengow/vngo/engine/backtesting_engine"
    "github.com/zhengow/vngo/gateway"
)

func NewBacktestingEngine() *backtesting_engine.BacktestingEngine {
    return backtesting_engine.NewBacktestingEngine()
}

func NewLiveTradeEngine(gI gateway.GatewayInterface) *engine.LiveTradeEngine {
    return engine.NewLiveTradeEngine(gI)
}
