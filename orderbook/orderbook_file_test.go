package orderbook

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func testOrderBookFromFile(t *testing.T, filename string) {
	orderBook := createOrderBook(t)

	file, err := os.Open(filename)
	require.NoError(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var expectedTotal, expectedBids, expectedAsks int

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "A":
			side, err := StringToOrderSide(fields[1])
			require.NoError(t, err)
			orderType, err := StringToOrderType(fields[2])
			require.NoError(t, err)
			price, err := strconv.Atoi(fields[3])
			require.NoError(t, err)
			qty, err := strconv.Atoi(fields[4])
			require.NoError(t, err)
			order := CreateOrder(orderType, side, Price(price), Quantity(qty))
			orderBook.AddOrder(order)
		case "R":
			expectedTotal, err = strconv.Atoi(fields[1])
			require.NoError(t, err)
			expectedBids, err = strconv.Atoi(fields[2])
			require.NoError(t, err)
			expectedAsks, err = strconv.Atoi(fields[3])
			require.NoError(t, err)
		default:
			t.Fatalf("unknown line prefix: %s", fields[0])
		}
	}

	require.NoError(t, scanner.Err())

	actualTotal := len(orderBook.Orders)
	actualBids := len(orderBook.Bids.Values())
	actualAsks := len(orderBook.Asks.Values())

	require.Equal(t, expectedTotal, actualTotal, "Total orders mismatch")
	require.Equal(t, expectedBids, actualBids, "Bid orders mismatch")
	require.Equal(t, expectedAsks, actualAsks, "Ask orders mismatch")
}

func TestCase1(t *testing.T) {
	testOrderBookFromFile(t, "testcases/case1.txt")
	testOrderBookFromFile(t, "testcases/case2.txt")
	testOrderBookFromFile(t, "testcases/case3.txt")
}
