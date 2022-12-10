package vngo

import (
    "database/sql/driver"
    "fmt"
    "time"
)

type Bar struct {
    Symbol
    Datetime     DatabaseTime `gorm:"type:datetime"`
    Interval     Interval
    Volume       float64
    OpenInterest float64
    OpenPrice    float64
    HighPrice    float64
    LowPrice     float64
    ClosePrice   float64
}

func (b *Bar) FullName() string {
    return fmt.Sprintf("%s.%s", b.Symbol, b.Exchange)
}

func (b *Bar) GetKLineData() [4]float64 {
    return [4]float64{b.OpenPrice, b.ClosePrice, b.LowPrice, b.HighPrice}
}

type Order struct {
    Symbol    Symbol
    OrderId   int
    Direction Direction
    Price     float64
    Volume    float64
    //status    Status
    //datetime time.Time
}

func NewOrder(symbol Symbol,
    orderId int,
    direction Direction,
    price float64,
    volume float64,
//status Status,
//    datetime time.Time
) *Order {
    return &Order{
        Symbol:    symbol,
        OrderId:   orderId,
        Direction: direction,
        Price:     price,
        Volume:    volume,
        //status:    status,
        //datetime: datetime,
    }
}

type TradeData struct {
    Symbol    Symbol
    OrderId   int
    TradeId   int
    Direction Direction
    Price     float64
    Volume    float64
    Datetime  time.Time
}

func NewTradeData(symbol Symbol,
    orderId int,
    tradeId int,
    direction Direction,
    price float64,
    volume float64,
    datetime time.Time,
) *TradeData {
    return &TradeData{
        Symbol:    symbol,
        OrderId:   orderId,
        TradeId:   tradeId,
        Direction: direction,
        Price:     price,
        Volume:    volume,
        Datetime:  datetime,
    }
}

func (t *TradeData) IsBuy() bool {
    return t.Direction == DirectionEnum.LONG
}

func (t *TradeData) IsSell() bool {
    return t.Direction == DirectionEnum.SHORT
}

type DatabaseTime struct {
    time.Time
}

func NewDatabaseTime(t time.Time) DatabaseTime {
    return DatabaseTime{
        t,
    }
}

func (t *DatabaseTime) UnmarshalJSON(data []byte) (err error) {
    now, err := time.ParseInLocation(DateFormat, string(data), time.Local)
    *t = DatabaseTime{
        now,
    }
    return
}

func (t *DatabaseTime) MarshalJSON() ([]byte, error) {
    tTime := t.Time
    return []byte(tTime.Format(DateFormat)), nil
}

func (t DatabaseTime) Value() (driver.Value, error) {
    return t.Format(DateFormat), nil
}

func (t *DatabaseTime) Scan(v interface{}) error {
    switch vt := v.(type) {
    case string:
        tTime, err := time.Parse(DateFormat, vt)
        if err != nil {
            return err
        }
        *t = DatabaseTime{
            tTime,
        }
    case time.Time:
        *t = DatabaseTime{
            vt,
        }
    default:
        return fmt.Errorf("unknown err: %v", v)
    }
    return nil
}
