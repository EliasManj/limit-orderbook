package orderbook

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func createOrderBook(t *testing.T) *OrderBook {
	orderbook := NewOrderBook()
	require.NotNil(t, orderbook)
	return orderbook
}

func CreateOrder(orderType orderType, side Side, price Price, initialQty Quantity) Order {
	orderId := RandomOrderId()
	return Order{
		orderId:      orderId,
		OrderType:    orderType,
		Side:         side,
		Price:        price,
		initialQty:   initialQty,
		remainingQty: initialQty,
	}
}

func TestOrderbook_AddSingleBids(t *testing.T) {
	orderbook := createOrderBook(t)
	// Generate 5 random price-amount pairs and test each case
	for i := 0; i < 5; i++ {
		price := RandomPrice()
		amount := RandomAmount()
		order := CreateOrder(GoodTilCancelled, Buy, price, amount)
		trades := orderbook.AddOrder(order)
		obOrder := orderbook.Orders[order.orderId]
		require.Equal(t, obOrder.orderId, order.orderId)
		require.Equal(t, obOrder.OrderType, order.OrderType)
		require.Equal(t, obOrder.Side, order.Side)
		require.Equal(t, obOrder.Price, order.Price)
		require.Len(t, trades, 0)
		require.Len(t, orderbook.Bids.Values(), i+1) // Bids should increase by 1 in each iteration
		bid := orderbook.Bids.Values()[price][0]
		require.Equal(t, bid.orderId, order.orderId)
		require.Equal(t, bid.OrderType, order.OrderType)
		require.Equal(t, bid.Side, order.Side)
		require.Equal(t, bid.Price, order.Price)
	}
}

func TestOrderbook_AddSingleAsks(t *testing.T) {
	orderbook := createOrderBook(t)
	// Generate 5 random price-amount pairs and test each case
	for i := 0; i < 5; i++ {
		price := RandomPrice()
		amount := RandomAmount()
		order := CreateOrder(GoodTilCancelled, Sell, price, amount)
		trades := orderbook.AddOrder(order)
		obOrder := orderbook.Orders[order.orderId]
		require.Equal(t, obOrder.orderId, order.orderId)
		require.Equal(t, obOrder.OrderType, order.OrderType)
		require.Equal(t, obOrder.Side, order.Side)
		require.Equal(t, obOrder.Price, order.Price)
		require.Len(t, trades, 0)
		require.Len(t, orderbook.Asks.Values(), i+1) // Bids should increase by 1 in each iteration
		bid := orderbook.Asks.Values()[price][0]
		require.Equal(t, bid.orderId, order.orderId)
		require.Equal(t, bid.OrderType, order.OrderType)
		require.Equal(t, bid.Side, order.Side)
		require.Equal(t, bid.Price, order.Price)
	}
}

func TestOrderbook_AddBidAndAskNoMatch(t *testing.T) {
	orderbook := createOrderBook(t)
	for i := 0; i < 5; i++ {
		price := RandomPrice()
		amount := RandomAmount()
		bid := CreateOrder(GoodTilCancelled, Buy, price, amount)
		ask := CreateOrder(GoodTilCancelled, Sell, price+100, amount)
		trades := orderbook.AddOrder(bid)
		require.Len(t, trades, 0)
		trades = orderbook.AddOrder(ask)
		require.Len(t, trades, 0)
		require.Len(t, orderbook.Bids.Values(), i+1)
		require.Len(t, orderbook.Asks.Values(), i+1)
	}
}

func TestOrderbook_BidAskMatch(t *testing.T) {
	orderbook := createOrderBook(t)
	price := RandomPrice()
	amount := RandomAmount()
	bid := CreateOrder(GoodTilCancelled, Buy, price, amount)
	ask := CreateOrder(GoodTilCancelled, Sell, price, amount)
	trades := orderbook.AddOrder(bid)
	require.Len(t, trades, 0)
	trades = orderbook.AddOrder(ask)
	require.Len(t, trades, 1)
	require.Len(t, orderbook.Bids.Values(), 0)
	require.Len(t, orderbook.Asks.Values(), 0)
	require.Equal(t, trades[0].BidTrade.Price, bid.Price)
	require.Equal(t, trades[0].AskTrade.Price, ask.Price)
	require.Equal(t, trades[0].BidTrade.Qty, bid.initialQty)
	require.Equal(t, trades[0].AskTrade.Qty, ask.initialQty)
	require.Equal(t, trades[0].BidTrade.OrderId, bid.orderId)
	require.Equal(t, trades[0].AskTrade.OrderId, ask.orderId)
}
