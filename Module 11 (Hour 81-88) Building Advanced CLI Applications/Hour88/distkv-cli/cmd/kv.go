package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Subhash-Pal/distkv-cli/internal/kvstore"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newPutCmd())
	rootCmd.AddCommand(newGetCmd())
	rootCmd.AddCommand(newDeleteCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newWatchCmd())
	rootCmd.AddCommand(newHealthCmd())
}

func newPutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "put <key> <value>",
		Short: "Store or update a key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newSignalContext(cfg.Timeout)
			defer cancel()

			return withClient(ctx, func(client *kvstore.Client) error {
				revision, err := client.Put(ctx, args[0], []byte(args[1]))
				if err != nil {
					return err
				}
				if cfg.JSON {
					return printJSON(map[string]any{
						"key":      args[0],
						"revision": revision,
					})
				}
				return printMessage("stored %q at revision %d", args[0], revision)
			})
		},
	}
}

func newGetCmd() *cobra.Command {
	var revision uint64

	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Fetch a key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newSignalContext(cfg.Timeout)
			defer cancel()

			return withClient(ctx, func(client *kvstore.Client) error {
				entry, err := client.Get(ctx, args[0], revision)
				if err != nil {
					return err
				}
				if cfg.JSON {
					return printJSON(map[string]any{
						"bucket":   entry.Bucket(),
						"key":      entry.Key(),
						"value":    string(entry.Value()),
						"revision": entry.Revision(),
						"created":  entry.Created().Format(time.RFC3339),
					})
				}
				return printMessage("%s=%s (rev=%d)", entry.Key(), string(entry.Value()), entry.Revision())
			})
		},
	}

	cmd.Flags().Uint64Var(&revision, "revision", 0, "Fetch a specific revision")
	return cmd
}

func newDeleteCmd() *cobra.Command {
	var purge bool

	cmd := &cobra.Command{
		Use:   "delete <key>",
		Short: "Delete a key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newSignalContext(cfg.Timeout)
			defer cancel()

			return withClient(ctx, func(client *kvstore.Client) error {
				if purge {
					if err := client.Purge(ctx, args[0]); err != nil {
						return err
					}
				} else {
					if err := client.Delete(ctx, args[0]); err != nil {
						return err
					}
				}
				return printMessage("removed %q", args[0])
			})
		},
	}

	cmd.Flags().BoolVar(&purge, "purge", false, "Remove all history for the key")
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List keys in the bucket",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newSignalContext(cfg.Timeout)
			defer cancel()

			return withClient(ctx, func(client *kvstore.Client) error {
				keys, err := client.List(ctx)
				if err != nil {
					return err
				}
				if cfg.JSON {
					return printJSON(map[string]any{"keys": keys})
				}
				if len(keys) == 0 {
					return printMessage("no keys found")
				}
				_, err = fmt.Fprintln(os.Stdout, strings.Join(keys, "\n"))
				return err
			})
		},
	}
}

func newWatchCmd() *cobra.Command {
	var updatesOnly bool

	cmd := &cobra.Command{
		Use:   "watch [pattern]",
		Short: "Watch key updates",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := ">"
			if len(args) == 1 {
				pattern = args[0]
			}

			ctx, cancel := newSignalContext(24 * time.Hour)
			defer cancel()

			return withClient(ctx, func(client *kvstore.Client) error {
				updates, err := client.Watch(ctx, pattern, updatesOnly)
				if err != nil {
					return err
				}
				for update := range updates {
					if update == nil {
						if !cfg.JSON {
							fmt.Fprintln(os.Stdout, "-- watcher caught up --")
						}
						continue
					}
					if cfg.JSON {
						if err := printJSON(map[string]any{
							"key":       update.Key(),
							"revision":  update.Revision(),
							"operation": normalizeWatchOp(update.Operation()),
							"value":     string(update.Value()),
						}); err != nil {
							return err
						}
						continue
					}
					fmt.Fprintf(os.Stdout, "%s rev=%d op=%s value=%q\n", update.Key(), update.Revision(), normalizeWatchOp(update.Operation()), string(update.Value()))
				}
				return nil
			})
		},
	}

	cmd.Flags().BoolVar(&updatesOnly, "updates-only", false, "Only stream new changes")
	return cmd
}

func newHealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check NATS and JetStream connectivity",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newSignalContext(cfg.Timeout)
			defer cancel()

			return withClient(ctx, func(client *kvstore.Client) error {
				report, err := client.Health(ctx)
				if err != nil {
					return err
				}
				if cfg.JSON {
					return printJSON(report)
				}
				return printMessage("server=%s jetstream=%t bucket=%s", report.Server, report.JetStream, report.Bucket)
			})
		},
	}
}
