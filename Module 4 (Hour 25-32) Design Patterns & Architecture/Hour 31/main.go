package main

import "fmt"

type Event struct {
	Name string
	Data string
}



type EventHandler func(Event)

type EventBus struct {
	subscribers map[string][]EventHandler//example: "OrderCreated" -> [inventoryHandler, emailHandler]
    //subscribers["OrderCreated"] = []EventHandler{inventoryHandler, emailHandler}
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

func (b *EventBus) Subscribe(eventName string, handler EventHandler) {
	b.subscribers[eventName] = append(b.subscribers[eventName], handler)
}

func (b *EventBus) Publish(event Event) {
	fmt.Println("publishing event:", event.Name)
	for _, handler := range b.subscribers[event.Name] {
		handler(event)
	}
}
/////////////
type OrderService struct {
	bus *EventBus
}

func NewOrderService(bus *EventBus) *OrderService {
	return &OrderService{bus: bus}
}

func (s *OrderService) CreateOrder(orderID string) {
	fmt.Println("order created:", orderID)
	s.bus.Publish(Event{
		Name: "OrderCreated",
		Data: orderID,
	})
}

func inventoryHandler(event Event) {
	fmt.Println("inventory reserved for order:", event.Data)
}

func emailHandler(event Event) {
	fmt.Println("email sent for order:", event.Data)
}

func main() {
	bus := NewEventBus()
	bus.Subscribe("OrderCreated", inventoryHandler)
	bus.Subscribe("OrderCreated", emailHandler)

	orderService := NewOrderService(bus)
	orderService.CreateOrder("ORD-1001")
}
