✅ 1. Go to project folder
cd "D:\training_golang\Module 10 (Hour 73-80) Event-Driven Distributed Systems with NATS & JetStream\order-system"

✅ 2. Start NATS (Terminal 1)
nats-server -js
Wait until:
Server is ready

✅ 3. Run Setup (Terminal 2)
go run cmd/setup/main.go

✅ 4. Start Processor (Terminal 3)
go run cmd/processor/main.go

✅ 5. Run Publisher (Terminal 4)
go run cmd/publisher/main.go

✅ 6. Check Status (Terminal 5 or reuse Terminal 2)
go run cmd/status/main.go

🧪 Optional (Crash Test)
# Stop processor (Ctrl + C), then restartgo run cmd/processor/main.go

🎯 Expected Final Output
order-1 → PROCESSEDorder-2 → PROCESSEDorder-3 → PROCESSED

