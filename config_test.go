package vngo

import (
    "github.com/zhengow/vngo/models"
    "testing"
)

func TestSymbol_Compare(t *testing.T) {
    s1 := models.Symbol{
        Name:     "1",
        Exchange: "1",
    }
    s2 := models.Symbol{
        Name:     "1",
        Exchange: "1",
    }

    println(s1 == s2)
}
