package kvstore

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Config struct {
	Server string
	Bucket string
}

type BucketOptions struct {
	History      int64
	TTL          time.Duration
	MaxBytes     int64
	MaxValueSize int32
}

type HealthReport struct {
	Server    string `json:"server"`
	JetStream bool   `json:"jetstream"`
	Bucket    string `json:"bucket"`
}

type Client struct {
	conn   *nats.Conn
	js     jetstream.JetStream
	bucket string
}

func Connect(ctx context.Context, cfg Config) (*Client, error) {
	nc, err := nats.Connect(cfg.Server, nats.Name("distkv-cli"))
	if err != nil {
		return nil, fmt.Errorf("connect to NATS: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("create jetstream client: %w", err)
	}

	return &Client{
		conn:   nc,
		js:     js,
		bucket: cfg.Bucket,
	}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) EnsureBucket(ctx context.Context, opts BucketOptions) (bool, jetstream.KeyValueStatus, error) {
	kv, err := c.js.KeyValue(ctx, c.bucket)
	if err == nil {
		status, statusErr := kv.Status(ctx)
		return false, status, statusErr
	}
	if !errors.Is(err, jetstream.ErrBucketNotFound) {
		return false, nil, err
	}

	kv, err = c.js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:       c.bucket,
		History:      normalizeHistory(opts.History),
		TTL:          opts.TTL,
		MaxBytes:     opts.MaxBytes,
		MaxValueSize: opts.MaxValueSize,
	})
	if err != nil {
		return false, nil, fmt.Errorf("create bucket: %w", err)
	}
	status, err := kv.Status(ctx)
	return true, status, err
}

func (c *Client) BucketStatus(ctx context.Context) (jetstream.KeyValueStatus, error) {
	kv, err := c.js.KeyValue(ctx, c.bucket)
	if err != nil {
		return nil, err
	}
	return kv.Status(ctx)
}

func (c *Client) Put(ctx context.Context, key string, value []byte) (uint64, error) {
	kv, err := c.js.KeyValue(ctx, c.bucket)
	if err != nil {
		return 0, err
	}
	return kv.Put(ctx, key, value)
}

func (c *Client) Get(ctx context.Context, key string, revision uint64) (jetstream.KeyValueEntry, error) {
	kv, err := c.js.KeyValue(ctx, c.bucket)
	if err != nil {
		return nil, err
	}
	if revision > 0 {
		return kv.GetRevision(ctx, key, revision)
	}
	return kv.Get(ctx, key)
}

func (c *Client) Delete(ctx context.Context, key string) error {
	kv, err := c.js.KeyValue(ctx, c.bucket)
	if err != nil {
		return err
	}
	return kv.Delete(ctx, key)
}

func (c *Client) Purge(ctx context.Context, key string) error {
	kv, err := c.js.KeyValue(ctx, c.bucket)
	if err != nil {
		return err
	}
	return kv.Purge(ctx, key)
}

func (c *Client) List(ctx context.Context) ([]string, error) {
	kv, err := c.js.KeyValue(ctx, c.bucket)
	if err != nil {
		return nil, err
	}
	lister, err := kv.ListKeys(ctx)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0)
	for key := range lister.Keys() {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys, nil
}

func (c *Client) Watch(ctx context.Context, pattern string, updatesOnly bool) (<-chan jetstream.KeyValueEntry, error) {
	kv, err := c.js.KeyValue(ctx, c.bucket)
	if err != nil {
		return nil, err
	}
	opts := make([]jetstream.WatchOpt, 0, 1)
	if updatesOnly {
		opts = append(opts, jetstream.UpdatesOnly())
	}
	watcher, err := kv.Watch(ctx, pattern, opts...)
	if err != nil {
		return nil, err
	}
	out := make(chan jetstream.KeyValueEntry)
	go func() {
		defer close(out)
		defer watcher.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case entry, ok := <-watcher.Updates():
				if !ok {
					return
				}
				out <- entry
			}
		}
	}()
	return out, nil
}

func (c *Client) Health(ctx context.Context) (HealthReport, error) {
	account, err := c.js.AccountInfo(ctx)
	if err != nil {
		return HealthReport{}, err
	}
	return HealthReport{
		Server:    c.conn.ConnectedUrl(),
		JetStream: account != nil,
		Bucket:    c.bucket,
	}, nil
}

func normalizeHistory(history int64) uint8 {
	if history <= 0 {
		return 1
	}
	if history > 64 {
		return 64
	}
	return uint8(history)
}
