I’ve created a complete **README.md in the canvas** with:

* Full setup
* Code structure
* Strict execution order
* Crash test procedure
* Retention policy experiments
* Debugging section

---

## Run Order (Quick Reference)

Follow this sequence every time:

```bash
# 1. Start NATS with JetStream
nats-server -js

# 2. Create stream (only once)
go run stream.go

# 3. Start durable consumer
go run consumer.go

# 4. Publish messages
go run publisher.go
```

---

## Mandatory Validation

After running:

1. Kill consumer during processing
2. Restart:

```bash
go run consumer.go
```

✔ You must see **re-delivery of unprocessed messages**

---
