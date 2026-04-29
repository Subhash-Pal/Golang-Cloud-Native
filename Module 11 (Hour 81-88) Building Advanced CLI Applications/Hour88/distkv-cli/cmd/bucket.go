package cmd

import (
	"fmt"
	"time"

	"github.com/Subhash-Pal/distkv-cli/internal/kvstore"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newBucketCmd())
}

func newBucketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bucket",
		Short: "Manage the JetStream key-value bucket",
	}

	cmd.AddCommand(newBucketCreateCmd())
	cmd.AddCommand(newBucketInfoCmd())
	return cmd
}

func newBucketCreateCmd() *cobra.Command {
	var history int64
	var ttl time.Duration
	var maxBytes int64
	var maxValue int32

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create the configured bucket if it does not already exist",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newSignalContext(cfg.Timeout)
			defer cancel()

			return withClient(ctx, func(client *kvstore.Client) error {
				created, status, err := client.EnsureBucket(ctx, kvstore.BucketOptions{
					History:      history,
					TTL:          ttl,
					MaxBytes:     maxBytes,
					MaxValueSize: maxValue,
				})
				if err != nil {
					return err
				}
				if cfg.JSON {
					return printJSON(map[string]any{
						"created": created,
						"bucket":  status.Bucket(),
						"history": status.History(),
						"ttl":     status.TTL().String(),
						"values":  status.Values(),
					})
				}
				action := "reused"
				if created {
					action = "created"
				}
				return printMessage("bucket %q %s", status.Bucket(), action)
			})
		},
	}

	cmd.Flags().Int64Var(&history, "history", 5, "Number of revisions to keep per key")
	cmd.Flags().DurationVar(&ttl, "ttl", 0, "Optional TTL for keys")
	cmd.Flags().Int64Var(&maxBytes, "max-bytes", 0, "Optional bucket max bytes")
	cmd.Flags().Int32Var(&maxValue, "max-value-size", 0, "Optional max value size in bytes")
	return cmd
}

func newBucketInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show bucket status",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newSignalContext(cfg.Timeout)
			defer cancel()

			return withClient(ctx, func(client *kvstore.Client) error {
				status, err := client.BucketStatus(ctx)
				if err != nil {
					return err
				}
				if cfg.JSON {
					return printJSON(map[string]any{
						"bucket":       status.Bucket(),
						"values":       status.Values(),
						"history":      status.History(),
						"ttl":          status.TTL().String(),
						"bytes":        status.Bytes(),
						"backingStore": fmt.Sprint(status.BackingStore()),
					})
				}
				return printMessage("bucket=%s values=%d history=%d ttl=%s bytes=%d", status.Bucket(), status.Values(), status.History(), status.TTL(), status.Bytes())
			})
		},
	}
}
