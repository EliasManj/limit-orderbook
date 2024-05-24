package orderbook

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOrderId() OrderId {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	timestamp = timestamp[len(timestamp)-10:]  // Extract the last 10 digits
	intTimestamp, _ := strconv.Atoi(timestamp) // Convert to int
	return OrderId(intTimestamp)
}

func RandomPrice() Price {
	// Generate a random price between 1 and 1000
	return Price(rand.Intn(1000) + 1)
}

func RandomAmount() Quantity {
	// Generate a random amount between 1 and 100
	return Quantity(rand.Intn(100) + 1)
}

// StringToOrderType converts a string to an orderType enumeration
func StringToOrderType(s string) (orderType, error) {
	switch strings.ToLower(s) {
	case "goodtillcancel", "goodtilcancelled":
		return GoodTilCancelled, nil
	case "fillandkill":
		return FillAndKill, nil
	case "fillorkill":
		return FillOrKill, nil
	case "market":
		return Market, nil
	default:
		return -1, errors.New("invalid order type: " + s)
	}
}

func StringToOrderSide(s string) (Side, error) {
	switch strings.ToLower(s) {
	case "b", "buy":
		return Buy, nil
	case "s", "sell":
		return Sell, nil
	default:
		return -1, errors.New("invalid order type: " + s)
	}
}
