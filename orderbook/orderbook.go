package orderbook

import (
	"errors"
)

// Enumerations and types
type orderType int
type Side int
type Quantity int32
type Price float64
type OrderId int

const (
	GoodTilCancelled orderType = iota
	FillAndKill
	FillOrKill
	Market
)

func (o orderType) String() string {
	switch o {
	case GoodTilCancelled:
		return "GoodTilCancelled"
	case FillAndKill:
		return "FillAndKill"
	case FillOrKill:
		return "FillOrKill"
	case Market:
		return "Market"
	default:
		return "Unknown"
	}
}

const (
	Buy Side = iota
	Sell
)

func (s Side) String() string {
	switch s {
	case Buy:
		return "Buy"
	case Sell:
		return "Sell"
	default:
		return "Unknown"
	}
}

// End Enumerations

// Orderbook definition start
type LevelInfo struct {
	Price    Price
	Quantity Quantity
}

type OrderBookLevelInfos struct {
	Bids []LevelInfo
	Asks []LevelInfo
}

// Orderbook definition end

// Order definition start

type Order struct {
	OrderType    orderType
	Side         Side
	orderId      OrderId
	Price        Price
	initialQty   Quantity
	remainingQty Quantity
}

func (o *Order) GetOrderId() OrderId {
	return o.orderId
}

func NewOrder(otype string, side string, price float64, qty int) *Order {
	orderType, err := StringToOrderType(otype)
	if err != nil {
		return nil
	}
	orderSide, err := StringToOrderSide(side)
	if err != nil {
		return nil
	}
	return &Order{
		OrderType:    orderType,
		Side:         orderSide,
		Price:        Price(price),
		initialQty:   Quantity(qty),
		remainingQty: Quantity(qty),
		orderId:      RandomOrderId(),
	}
}

func (o *Order) GetInitialQty() Quantity {
	return o.initialQty
}

func (o *Order) GetFilledQty() Quantity {
	return o.initialQty - o.remainingQty
}

func (o *Order) Fill(qty Quantity) error {
	if qty > o.remainingQty {
		return errors.New("cannot fill more than remaining quantity")
	}
	o.remainingQty -= qty
	return nil
}

func (o *Order) IsFilled() bool {
	return o.remainingQty == 0
}

type TradeInfo struct {
	OrderId OrderId
	Price   Price
	Qty     Quantity
}

type Trade struct {
	BidTrade TradeInfo
	AskTrade TradeInfo
}

var Trades = []Trade{}
var Orders = []Trade{}

// Order definition end
type OrderBook struct {
	Bids   *OrderedMap
	Asks   *OrderedMap
	Orders map[OrderId]Order
}

// NewOrderBook creates a new OrderBook with Bids in ascending order and Asks in descending order
func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids:   NewOrderedMap(Descending),
		Asks:   NewOrderedMap(Ascending),
		Orders: make(map[OrderId]Order),
	}
}

func (ob *OrderBook) CanMatch(side Side, price Price) bool {
	if side == Buy {
		if len(ob.Asks.Keys()) == 0 {
			return false
		}
		bestAskPrice, _ := ob.Asks.FirstKey()
		return price >= bestAskPrice
	} else {
		if len(ob.Bids.Keys()) == 0 {
			return false
		}
		bestBidPrice, _ := ob.Bids.FirstKey()
		return price <= bestBidPrice
	}
}

func (ob *OrderBook) GetTotalQty(side Side, price Price) Quantity {
	var sum Quantity = 0
	if side == Buy {
		list := ob.Bids.Values()[price]
		for _, value := range list {
			sum += value.remainingQty
		}
	} else {
		list := ob.Asks.Values()[price]
		for _, value := range list {
			sum += value.remainingQty
		}

	}
	return sum
}

func (ob *OrderBook) CanMatchCompletely(side Side, price Price, quantity Quantity) bool {
	if side == Buy {
		if len(ob.Asks.Keys()) == 0 {
			return false
		}
		bestAskPrice, _ := ob.Asks.FirstKey()
		maxQty := ob.GetTotalQty(Sell, price)
		return price >= bestAskPrice && quantity <= maxQty
	} else {
		if len(ob.Bids.Keys()) == 0 {
			return false
		}
		bestBidPrice, _ := ob.Bids.FirstKey()
		maxQty := ob.GetTotalQty(Buy, price)
		return price <= bestBidPrice && quantity <= maxQty
	}
}

func (ob *OrderBook) MatchOrders() []Trade {
	trades := []Trade{}
	for {
		if ob.Bids.IsEmpty() || ob.Asks.IsEmpty() {
			break
		}
		bidPrice, bids := ob.Bids.BestPrice()
		askPrice, asks := ob.Asks.BestPrice()
		if bidPrice < askPrice {
			break
		}
		for len(bids) > 0 && len(asks) > 0 {
			bid := &bids[0]
			ask := &asks[0]
			quantity := min(bid.remainingQty, ask.remainingQty)
			bid.Fill(quantity)
			ask.Fill(quantity)
			if bid.IsFilled() {
				ob.Bids.Values()[bidPrice] = ob.Bids.Values()[bidPrice][1:]
				bids = bids[1:]
				delete(ob.Orders, bid.orderId)
			}
			if ask.IsFilled() {
				ob.Asks.Values()[askPrice] = ob.Asks.Values()[askPrice][1:]
				asks = asks[1:]
				delete(ob.Orders, ask.orderId)
			}
			if len(ob.Bids.Values()[bidPrice]) == 0 {
				ob.Bids.Delete(bidPrice)
			}
			if len(ob.Asks.Values()[askPrice]) == 0 {
				ob.Asks.Delete(askPrice)
			}
			trades = append(trades, Trade{
				BidTrade: TradeInfo{
					OrderId: bid.orderId,
					Price:   bid.Price,
					Qty:     quantity,
				},
				AskTrade: TradeInfo{
					OrderId: ask.orderId,
					Price:   ask.Price,
					Qty:     quantity,
				},
			})
		}
	}
	if !ob.Bids.IsEmpty() {
		_, bids := ob.Bids.BestPrice()
		order := bids[0]
		if order.OrderType == FillAndKill {
			ob.CancelOrder(order.orderId)
		}
	}
	if !ob.Asks.IsEmpty() {
		_, asks := ob.Asks.BestPrice()
		order := asks[0]
		if order.OrderType == FillAndKill {
			ob.CancelOrder(order.orderId)
		}
	}
	return trades
}

func (ob *OrderBook) AddOrder(order Order) []Trade {
	if ob.Orders[order.orderId] != (Order{}) {
		return nil
	}
	if order.OrderType == Market {
		if order.Side == Buy && !ob.Asks.IsEmpty() {
			worstAskPrice, _ := ob.Asks.LastKey()
			order.Price = worstAskPrice
			order.OrderType = GoodTilCancelled
		} else if order.Side == Sell && !ob.Bids.IsEmpty() {
			worstBidPrice, _ := ob.Bids.LastKey()
			order.Price = worstBidPrice
			order.OrderType = GoodTilCancelled
		} else {
			return nil
		}
	}
	if order.OrderType == FillAndKill && !ob.CanMatch(order.Side, order.Price) {
		return nil
	}
	if order.OrderType == FillOrKill && !ob.CanMatchCompletely(order.Side, order.Price, order.initialQty) {
		return nil
	}
	if order.Side == Buy {
		ob.Bids.Add(order.Price, order)
	} else {
		ob.Asks.Add(order.Price, order)
	}
	ob.Orders[order.orderId] = order
	return ob.MatchOrders()
}

func (ob *OrderBook) CancelOrder(orderId OrderId) {
	order := ob.Orders[orderId]
	if order == (Order{}) {
		return
	}
	if order.Side == Buy {
		ob.Bids.DeleteOrder(order)
	} else {
		ob.Asks.DeleteOrder(order)
	}
}

func (ob *OrderBook) ModifyOrder(order Order) {
	ob.CancelOrder(order.orderId)
	ob.AddOrder(order)
}

func (ob *OrderBook) Size() int {
	return len(ob.Orders)
}

func (ob *OrderBook) GetOrderInfos() OrderBookLevelInfos {
	bids := []LevelInfo{}
	for _, orders := range ob.Bids.Values() {
		for _, order := range orders {
			bids = append(bids, LevelInfo{
				Price:    order.Price,
				Quantity: order.remainingQty,
			})
		}
	}
	asks := []LevelInfo{}
	for _, orders := range ob.Asks.Values() {
		for _, order := range orders {
			asks = append(asks, LevelInfo{
				Price:    order.Price,
				Quantity: order.remainingQty,
			})
		}
	}
	return OrderBookLevelInfos{
		Bids: bids,
		Asks: asks,
	}
}
