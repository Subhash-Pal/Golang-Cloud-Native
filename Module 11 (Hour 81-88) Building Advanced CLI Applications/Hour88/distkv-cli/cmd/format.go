package cmd

import "github.com/nats-io/nats.go/jetstream"

func normalizeWatchOp(op jetstream.KeyValueOp) string {
	if op == 0 {
		return "put"
	}
	return op.String()
}
