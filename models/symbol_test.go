package models

import (
    "testing"
)

func TestSymbol_Compare(t *testing.T) {
    s1 := Symbol{
        Name:     "1",
        Exchange: "1",
    }
    s2 := Symbol{
        Name:     "1",
        Exchange: "1",
    }

    println(s1 == s2)
}
