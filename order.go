package vngo

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
