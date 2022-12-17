package models

type Order struct {
    Symbol    Symbol
    OrderId   string
    Direction Direction
    Price     float64
    Volume    float64
    Status    Status
    Traded    float64
    //datetime models.VnTime
}

func NewOrder(symbol Symbol,
    orderId string,
    direction Direction,
    price float64,
    volume float64,
//status Status,
//    datetime models.VnTime
) *Order {
    return &Order{
        Symbol:    symbol,
        OrderId:   orderId,
        Direction: direction,
        Price:     price,
        Volume:    volume,
        Status:    StatusEnum.SUBMITTING,
        //datetime: datetime,
    }
}
