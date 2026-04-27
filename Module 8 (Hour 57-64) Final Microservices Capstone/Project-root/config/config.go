package config

type ServiceConfig struct {
	AuthService    string
	ProductService string
	OrderService   string
}

func LoadConfig() ServiceConfig {
	return ServiceConfig{
		AuthService:    "http://localhost:8001",
		ProductService: "http://localhost:8002",
		OrderService:   "http://localhost:8003",
	}
}