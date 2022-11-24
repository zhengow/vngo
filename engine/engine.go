package engine

type BacktestingEngine struct {
}

func (b *BacktestingEngine) SetParameters () {
	println("set params")
}

func (b *BacktestingEngine) AddStrategy () {
	println("add")
}

func (b *BacktestingEngine) LoadData () {
	println("load")
}

func (b *BacktestingEngine) RunBacktesting () {
	println("run")
}

func (b *BacktestingEngine) CalculateResult () {
	println("calc")
}