package orderbook

import (
	"sort"
)

type MapOrder int

const (
	Ascending MapOrder = iota
	Descending
)

type OrderedMap struct {
	keys   []Price
	values map[Price][]Order
	order  MapOrder
}

// NewOrderedMap creates a new OrderedMap with the specified order
func NewOrderedMap(order MapOrder) *OrderedMap {
	return &OrderedMap{
		keys:   make([]Price, 0),
		values: make(map[Price][]Order),
		order:  order,
	}
}

// Add adds a key-value pair to the map and maintains the sorted order of keys
func (om *OrderedMap) Add(key Price, value Order) {
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
		om.sortKeys()
	}
	om.values[key] = append(om.values[key], value)
}

// Get retrieves a value by key
func (om *OrderedMap) Get(key Price) ([]Order, bool) {
	value, exists := om.values[key]
	return value, exists
}

// Delete removes a key-value pair from the map and updates the sorted order of keys
func (om *OrderedMap) Delete(key Price) {
	if _, exists := om.values[key]; exists {
		delete(om.values, key)
		for i, k := range om.keys {
			if k == key {
				om.keys = append(om.keys[:i], om.keys[i+1:]...)
				break
			}
		}
		om.sortKeys()
	}
}

// DeleteOrder removes a specific order from the order slice of the specified price
func (om *OrderedMap) DeleteOrder(order Order) {
	key := order.Price
	if orders, exists := om.values[key]; exists {
		for i, o := range orders {
			if o.orderId == order.orderId {
				om.values[key] = append(om.values[key][:i], om.values[key][i+1:]...)
				break
			}
		}
		if len(om.values[key]) == 0 {
			om.Delete(key)
		}
	}
}

// Keys returns the keys in the map in sorted order
func (om *OrderedMap) Keys() []Price {
	return om.keys
}

// Values returns the values in the map in sorted order of keys
func (om *OrderedMap) Values() map[Price][]Order {
	return om.values
}

// sortKeys sorts the keys according to the specified order
func (om *OrderedMap) sortKeys() {
	if om.order == Ascending {
		sort.Slice(om.keys, func(i, j int) bool { return om.keys[i] < om.keys[j] })
	} else if om.order == Descending {
		sort.Slice(om.keys, func(i, j int) bool { return om.keys[i] > om.keys[j] })
	}
}

// FirstKey returns the first key in the sorted order
func (om *OrderedMap) FirstKey() (Price, bool) {
	if len(om.keys) == 0 {
		var zero Price
		return zero, false
	}
	return om.keys[0], true
}

// FirstValue returns the first value in the sorted order
func (om *OrderedMap) FirstValue() ([]Order, bool) {
	if len(om.keys) == 0 {
		return nil, false
	}
	firstKey := om.keys[0]
	return om.values[firstKey], true
}

// FirstPrice returns the first key in the sorted order
func (om *OrderedMap) BestPrice() (Price, []Order) {
	if len(om.keys) == 0 {
		var zero Price
		return zero, nil
	}
	firstKey := om.keys[0]
	return firstKey, om.values[firstKey]
}

// LastKey returns the first key in the sorted order
func (om *OrderedMap) LastKey() (Price, bool) {
	if len(om.keys) == 0 {
		var zero Price
		return zero, false
	}
	return om.keys[len(om.keys)-1], true
}

// LastValue returns the last value in the sorted order
func (om *OrderedMap) LastValue() ([]Order, bool) {
	if len(om.keys) == 0 {
		return nil, false
	}
	lastKey := om.keys[len(om.keys)-1]
	return om.values[lastKey], true
}

// WorstPrice returns the last key in the sorted order
func (om *OrderedMap) LastPrice() (Price, []Order) {
	if len(om.keys) == 0 {
		var zero Price
		return zero, nil
	}
	lastKey := om.keys[len(om.keys)-1]
	return lastKey, om.values[lastKey]
}

// IsEmpty checks if the OrderedMap is empty
func (om *OrderedMap) IsEmpty() bool {
	return len(om.keys) == 0
}
