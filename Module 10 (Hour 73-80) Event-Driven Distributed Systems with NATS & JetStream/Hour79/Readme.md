Your **KV Store README.md** is now ready in the canvas.

---

## Quick Run Order (Checklist)

```bash
# 1. Start server
nats-server -js

# 2. Create KV bucket
go run kv_create.go

# 3. Start watcher (Terminal 1)
go run kv_watch.go

# 4. Put value (Terminal 2)
go run kv_put.go

# 5. Read value
go run kv_get.go
```

---

## Critical Validation

* Watcher must print updates immediately
* Revision should increment on each update
* Restarting watcher should still reflect latest state

---

