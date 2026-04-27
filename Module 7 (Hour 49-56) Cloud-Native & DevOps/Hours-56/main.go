package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

type item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

type order struct {
	ID       int    `json:"id"`
	ItemID   int    `json:"item_id"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

type store struct {
	mu     sync.Mutex
	items  []item
	orders []order
}

func newStore() *store {
	return &store{
		items: []item{
			{ID: 1, Name: "Laptop", Stock: 5},
			{ID: 2, Name: "Keyboard", Stock: 10},
		},
	}
}

func (s *store) listItems(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	writeJSON(w, http.StatusOK, s.items)
}

func (s *store) createItem(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Stock int    `json:"stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
		return
	}
	if input.Name == "" || input.Stock < 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name and non-negative stock are required"})
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	newItem := item{
		ID:    len(s.items) + 1,
		Name:  input.Name,
		Stock: input.Stock,
	}
	s.items = append(s.items, newItem)
	writeJSON(w, http.StatusCreated, newItem)
}

func (s *store) listOrders(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	writeJSON(w, http.StatusOK, s.orders)
}

func (s *store) createOrder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ItemID   int `json:"item_id"`
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
		return
	}
	if input.Quantity <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "quantity must be greater than zero"})
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.items {
		if s.items[i].ID != input.ItemID {
			continue
		}
		if s.items[i].Stock < input.Quantity {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "insufficient stock"})
			return
		}

		s.items[i].Stock -= input.Quantity
		newOrder := order{
			ID:       len(s.orders) + 1,
			ItemID:   input.ItemID,
			Quantity: input.Quantity,
			Status:   "created",
		}
		s.orders = append(s.orders, newOrder)
		writeJSON(w, http.StatusCreated, newOrder)
		return
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "item not found"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func newMux(db *store) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"hour":  56,
			"topic": "Deploy containerized API",
			"routes": []string{
				"GET /items",
				"POST /items",
				"GET /orders",
				"POST /orders",
				"GET /healthz",
				"GET /readyz",
			},
		})
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})
	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			db.listItems(w, r)
		case http.MethodPost:
			db.createItem(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			db.listOrders(w, r)
		case http.MethodPost:
			db.createOrder(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("hour 56 api listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, newMux(newStore())))
}
