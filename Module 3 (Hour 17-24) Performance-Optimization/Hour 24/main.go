package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type user struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type apiServer struct {
	users       []user
	cachedUsers []byte
}

func main() {
	server := newAPIServer()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	mux := http.NewServeMux()
	mux.HandleFunc("/slow", server.slowHandler)
	mux.HandleFunc("/optimized", server.optimizedHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("/debug/pprof/", http.DefaultServeMux)

	slog.Info("server started", "addr", "http://127.0.0.1:"+port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server failed", "error", err)
	}
}

func newAPIServer() *apiServer {
	users := make([]user, 15_000)
	for i := range users {
		users[i] = user{
			ID:    i + 1,
			Name:  "user-" + strconv.Itoa(i+1),
			Score: 15_000 - i,
		}
	}

	sorted := slices.Clone(users)
	slices.SortFunc(sorted, func(a, b user) int {
		return b.Score - a.Score
	})

	cachedUsers, _ := json.Marshal(sorted[:200])

	return &apiServer{
		users:       users,
		cachedUsers: cachedUsers,
	}
}

func (s *apiServer) slowHandler(w http.ResponseWriter, r *http.Request) {
	filter := strings.ToLower(r.URL.Query().Get("q"))

	cloned := slices.Clone(s.users)
	slices.SortFunc(cloned, func(a, b user) int {
		return b.Score - a.Score
	})

	var response []user
	for _, item := range cloned {
		if filter == "" || strings.Contains(strings.ToLower(item.Name), filter) {
			time.Sleep(50 * time.Microsecond)
			response = append(response, item)
		}
		if len(response) == 200 {
			break
		}
	}

	writeJSON(w, response)
}

func (s *apiServer) optimizedHandler(w http.ResponseWriter, r *http.Request) {
	filter := strings.ToLower(r.URL.Query().Get("q"))
	if filter == "" {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(s.cachedUsers)
		return
	}

	var response []user
	response = make([]user, 0, 200)
	for _, item := range s.users {
		if strings.Contains(strings.ToLower(item.Name), filter) {
			response = append(response, item)
		}
		if len(response) == 200 {
			break
		}
	}

	writeJSON(w, response)
}

var jsonPool = sync.Pool{
	New: func() any {
		buffer := make([]byte, 0, 32*1024)
		return &buffer
	},
}

func writeJSON(w http.ResponseWriter, value any) {
	bufferPtr := jsonPool.Get().(*[]byte)
	defer jsonPool.Put(bufferPtr)

	buffer := (*bufferPtr)[:0]
	data, err := json.Marshal(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buffer = append(buffer, data...)
	*bufferPtr = buffer

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(buffer)
}
