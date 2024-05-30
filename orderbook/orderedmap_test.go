package orderbook

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestOrderedMapBestPrices(t *testing.T) {
	bids := NewOrderedMap(Descending)
	bids.Add(100.0, Order{Price: 100.0})
	bids.Add(200.0, Order{Price: 200.0})
	bids.Add(150.0, Order{Price: 150.0})
	bids.Add(50.0, Order{Price: 50.0})
	bids.Add(250.0, Order{Price: 250.0})
	bestBid, _ := bids.BestPrice()
	require.Equal(t, Price(250.0), bestBid)
	asks := NewOrderedMap(Ascending)
	asks.Add(100.0, Order{Price: 100.0})
	asks.Add(200.0, Order{Price: 200.0})
	asks.Add(150.0, Order{Price: 150.0})
	asks.Add(50.0, Order{Price: 50.0})
	asks.Add(250.0, Order{Price: 250.0})
	bestAsk, _ := asks.BestPrice()
	require.Equal(t, Price(50.0), bestAsk)
}

func TestOrderedMapAddOrders(t *testing.T) {
	om := NewOrderedMap(Ascending)
	om.Add(100.0, Order{Price: 100.0})
	om.Add(200.0, Order{Price: 200.0})
	om.Add(150.0, Order{Price: 150.0})
	om.Add(50.0, Order{Price: 50.0})
	om.Add(250.0, Order{Price: 250.0})

	if len(om.keys) != 5 {
		t.Errorf("Expected 5 keys, got %d", len(om.keys))
	}
	if len(om.values) != 5 {
		t.Errorf("Expected 5 values, got %d", len(om.values))
	}
	if om.keys[0] != 50.0 {
		t.Errorf("Expected 50.0, got %f", om.keys[0])
	}
	if om.keys[1] != 100.0 {
		t.Errorf("Expected 100.0, got %f", om.keys[1])
	}
	if om.keys[2] != 150.0 {
		t.Errorf("Expected 150.0, got %f", om.keys[2])
	}
	if om.keys[3] != 200.0 {
		t.Errorf("Expected 200.0, got %f", om.keys[3])
	}
	if om.keys[4] != 250.0 {
		t.Errorf("Expected 250.0, got %f", om.keys[4])
	}
}

func TestOrderedMapRemoveDeleteOrders(t *testing.T) {
	om := NewOrderedMap(Ascending)
	om.Add(100.0, Order{orderId: 1, Price: 100.0})
	om.Add(100.0, Order{orderId: 2, Price: 100.0})
	om.Add(100.0, Order{orderId: 3, Price: 100.0})
	om.Add(100.0, Order{orderId: 4, Price: 100.0})
	om.Add(300.0, Order{orderId: 5, Price: 300.0})
	om.DeleteOrder(Order{orderId: 2, Price: 100.0})
	om.DeleteOrder(Order{orderId: 3, Price: 100.0})
	require.Equal(t, 2, len(om.values[100.0]))
	require.Equal(t, 1, len(om.values[300.0]))
	require.NotContains(t, om.values[100.0], Order{orderId: 2, Price: 100.0})
	require.NotContains(t, om.values[100.0], Order{orderId: 3, Price: 100.0})
	require.Contains(t, om.values[100.0], Order{orderId: 1, Price: 100.0})
	require.Contains(t, om.values[100.0], Order{orderId: 4, Price: 100.0})
}

func TestOrderedMapDeleteKeys(t *testing.T) {
	om := NewOrderedMap(Ascending)
	om.Add(100.0, Order{Price: 100.0})
	om.Add(200.0, Order{Price: 200.0})
	om.Add(200.0, Order{Price: 200.0})
	om.Add(150.0, Order{Price: 150.0})
	om.Add(50.0, Order{Price: 50.0})
	om.Delete(200.0)
	require.Equal(t, 3, len(om.keys))
	require.NotContains(t, om.keys, 200.0)
}

func TestOrderedMapFirstValues(t *testing.T) {
	om := NewOrderedMap(Ascending)
	om.Add(100.0, Order{Price: 100.0})
	om.Add(200.0, Order{Price: 200.0})
	om.Add(300.0, Order{Price: 300.0})
	om.Add(400.0, Order{Price: 400.0})
	om.Add(500.0, Order{Price: 500.0})
	key, _ := om.FirstKey()
	require.Equal(t, Price(100.0), key)
	om = NewOrderedMap(Descending)
	om.Add(100.0, Order{Price: 100.0})
	om.Add(200.0, Order{Price: 200.0})
	om.Add(300.0, Order{Price: 300.0})
	om.Add(400.0, Order{Price: 400.0})
	om.Add(500.0, Order{Price: 500.0})
	key, _ = om.FirstKey()
	require.Equal(t, Price(500.0), key)
}

func TestOrderedmapIsEmpty(t *testing.T) {
	om := NewOrderedMap(Ascending)
	require.True(t, om.IsEmpty())
	om.Add(100.0, Order{Price: 100.0})
	require.False(t, om.IsEmpty())
}

func TestOrderedmapLastValues(t *testing.T) {
	om := NewOrderedMap(Ascending)
	om.Add(100.0, Order{Price: 100.0})
	om.Add(200.0, Order{Price: 200.0})
	om.Add(300.0, Order{Price: 300.0})
	om.Add(400.0, Order{Price: 400.0})
	om.Add(500.0, Order{Price: 500.0})
	worstPrice, _ := om.LastKey()
	require.Equal(t, Price(500.0), worstPrice)
	om = NewOrderedMap(Descending)
	om.Add(100.0, Order{Price: 100.0})
	om.Add(200.0, Order{Price: 200.0})
	om.Add(300.0, Order{Price: 300.0})
	om.Add(400.0, Order{Price: 400.0})
	om.Add(500.0, Order{Price: 500.0})
	worstPrice, _ = om.LastKey()
	require.Equal(t, Price(100.0), worstPrice)
}
