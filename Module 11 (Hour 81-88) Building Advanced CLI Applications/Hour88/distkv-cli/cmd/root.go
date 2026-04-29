package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Subhash-Pal/distkv-cli/internal/kvstore"
	"github.com/spf13/cobra"
)

type Config struct {
	Server  string
	Bucket  string
	Timeout time.Duration
	JSON    bool
}

func (c Config) Validate() error {
	if c.Server == "" {
		return errors.New("server cannot be empty")
	}
	if c.Bucket == "" {
		return errors.New("bucket cannot be empty")
	}
	if c.Timeout <= 0 {
		return errors.New("timeout must be greater than zero")
	}
	return nil
}

var cfg = Config{}

var rootCmd = &cobra.Command{
	Use:   "distkv",
	Short: "Distributed key-value CLI backed by NATS JetStream",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return cfg.Validate()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfg.Server, "server", "nats://127.0.0.1:4222", "NATS server URL")
	rootCmd.PersistentFlags().StringVar(&cfg.Bucket, "bucket", "distkv", "JetStream key-value bucket")
	rootCmd.PersistentFlags().DurationVar(&cfg.Timeout, "timeout", 5*time.Second, "Request timeout")
	rootCmd.PersistentFlags().BoolVar(&cfg.JSON, "json", false, "Print structured JSON output")
}

func newSignalContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	base, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithTimeout(base, timeout)
	return ctx, func() {
		cancel()
		stop()
	}
}

func withClient(ctx context.Context, fn func(*kvstore.Client) error) error {
	client, err := kvstore.Connect(ctx, kvstore.Config{
		Server: cfg.Server,
		Bucket: cfg.Bucket,
	})
	if err != nil {
		return err
	}
	defer client.Close()
	return fn(client)
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func printMessage(format string, args ...any) error {
	if cfg.JSON {
		return printJSON(map[string]any{
			"message": fmt.Sprintf(format, args...),
		})
	}
	_, err := fmt.Fprintf(os.Stdout, format+"\n", args...)
	return err
}
