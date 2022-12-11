package vngo

import (
    "github.com/zhengow/vngo/strategy"
    "testing"
)

func TestSymbol_Compare(t *testing.T) {
    s1 := strategy.Symbol{
        Name:     "1",
        Exchange: "1",
    }
    s2 := strategy.Symbol{
        Name:     "1",
        Exchange: "1",
    }

    println(s1 == s2)
}
