package vngo

import (
    "github.com/zhengow/vngo/engine/backtesting"
    "github.com/zhengow/vngo/engine/live_trade"
    "github.com/zhengow/vngo/gateway"
)

func NewBacktestingEngine() *backtesting.Engine {
    return backtesting.NewBacktestingEngine()
}

func NewLiveTradeEngine(gI gateway.GatewayInterface) *live_trade.Engine {
    return live_trade.NewLiveTradeEngine(gI)
}
