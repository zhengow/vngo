package engine

import "github.com/zhengow/vngo/model"

type positionEngine struct {
	positions map[model.Symbol]float64
}

func newPositionEngine() *positionEngine {
	return &positionEngine{
		positions: make(map[model.Symbol]float64),
	}
}
func (p *positionEngine) GetPositions() map[model.Symbol]float64 {
	return p.positions
}

func (p *positionEngine) UpdatePositions(symbol model.Symbol, incrementPos float64) {
	p.positions[symbol] += incrementPos
}
