package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	port := getenv("PORT", "8080")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"hour":  52,
			"topic": "Kubernetes architecture",
			"cluster": map[string]string{
				"cluster_name": getenv("CLUSTER_NAME", "training-cluster"),
				"node_name":    getenv("NODE_NAME", "worker-1"),
				"namespace":    getenv("NAMESPACE", "default"),
				"pod_name":     getenv("POD_NAME", "architecture-demo"),
			},
			"core_components": []string{
				"API Server",
				"etcd",
				"Scheduler",
				"Controller Manager",
				"Kubelet",
				"Kube Proxy",
				"Pods and Services",
			},
		})
	})
	mux.HandleFunc("/components", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"control_plane": "API Server, Scheduler, Controller Manager, etcd",
			"worker_node":   "Kubelet, Kube Proxy, Pods",
			"service_role":  "Provides stable networking to a changing set of Pods",
		})
	})

	addr := ":" + port
	log.Printf("hour 52 architecture app listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
