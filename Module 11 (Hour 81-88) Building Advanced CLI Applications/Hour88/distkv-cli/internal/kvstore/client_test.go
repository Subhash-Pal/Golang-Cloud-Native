package kvstore

import "testing"

func TestNormalizeHistory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   int64
		want uint8
	}{
		{-1, 1},
		{0, 1},
		{1, 1},
		{5, 5},
		{100, 64},
	}

	for _, tt := range tests {
		if got := normalizeHistory(tt.in); got != tt.want {
			t.Fatalf("normalizeHistory(%d)=%d want %d", tt.in, got, tt.want)
		}
	}
}
