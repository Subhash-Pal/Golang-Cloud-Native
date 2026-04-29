package cmd

import (
	"testing"
	"time"
)

func TestConfigValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid",
			cfg: Config{
				Server:  "nats://127.0.0.1:4222",
				Bucket:  "distkv",
				Timeout: 5 * time.Second,
			},
		},
		{
			name: "missing server",
			cfg: Config{
				Bucket:  "distkv",
				Timeout: 5 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "missing bucket",
			cfg: Config{
				Server:  "nats://127.0.0.1:4222",
				Timeout: 5 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "non positive timeout",
			cfg: Config{
				Server: "nats://127.0.0.1:4222",
				Bucket: "distkv",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.cfg.Validate()
			if tt.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
