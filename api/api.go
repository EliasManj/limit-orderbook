// api.go
package api

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/EliasManj/orderbook/orderbook"
)

var ob *orderbook.OrderBook = orderbook.NewOrderBook()

type createOrderJson struct {
	OrderType string  `json:"order_type" binding:"required"`
	Side      string  `json:"side" binding:"required"`
	Price     float64 `json:"price"`
	Qty       int     `json:"qty" binding:"required"`
	OrderId   int     `json:"order_id"`
}

type createOrderResponse struct {
	Trades []orderbook.Trade `json:"trades"`
	Order  createOrderJson   `json:"order"`
}

func executeOrder(order orderbook.Order) createOrderResponse {
	trades := ob.AddOrder(order)
	orderResponse := createOrderJson{
		OrderType: order.OrderType.String(),
		Side:      order.Side.String(),
		Price:     float64(order.Price),
		Qty:       int(order.GetInitialQty()),
		OrderId:   int(order.GetOrderId()),
	}
	return createOrderResponse{
		Trades: trades,
		Order:  orderResponse,
	}
}

func GetBids(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bids := ob.GetOrderInfos().Bids
	json.NewEncoder(w).Encode(bids)
}

func GetAsks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	asks := ob.GetOrderInfos().Asks
	json.NewEncoder(w).Encode(asks)
}

// SeedOrders populates the orderbook with random bids and asks around a mid-price so the
// UI has data to display immediately on startup.
func SeedOrders() {
	const (
		midPrice  = 100.0
		tickSize  = 0.5
		numLevels = 10
	)
	for i := 1; i <= numLevels; i++ {
		bidPrice := midPrice - float64(i)*tickSize
		askPrice := midPrice + float64(i)*tickSize
		bidQty := rand.Intn(91) + 10
		askQty := rand.Intn(91) + 10

		bid := orderbook.NewOrder("GoodTilCancelled", "Buy", bidPrice, bidQty)
		ask := orderbook.NewOrder("GoodTilCancelled", "Sell", askPrice, askQty)
		ob.AddOrder(*bid)
		ob.AddOrder(*ask)
	}
}

func ClearOrders(w http.ResponseWriter, r *http.Request) {
	ob = orderbook.NewOrderBook()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "cleared"})
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req createOrderJson
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order := orderbook.NewOrder(req.OrderType, req.Side, req.Price, req.Qty)
	orderResonse := executeOrder(*order)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(orderResonse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
