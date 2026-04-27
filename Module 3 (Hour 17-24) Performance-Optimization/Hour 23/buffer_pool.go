package pooldemo

import (
	"bytes"
	"encoding/json"
	"sync"
)

type order struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Amount int    `json:"amount"`
}

var jsonBufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func EncodeWithoutPool(items []order) ([]byte, error) {
	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(items); err != nil {
		return nil, err
	}
	return append([]byte(nil), buffer.Bytes()...), nil
}

func EncodeWithPool(items []order) ([]byte, error) {
	buffer := jsonBufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer jsonBufferPool.Put(buffer)

	if err := json.NewEncoder(buffer).Encode(items); err != nil {
		return nil, err
	}
	return append([]byte(nil), buffer.Bytes()...), nil
}

func sampleOrders(size int) []order {
	items := make([]order, size)
	for i := range items {
		items[i] = order{
			ID:     i + 1,
			Status: "processed",
			Amount: (i + 1) * 100,
		}
	}
	return items
}
