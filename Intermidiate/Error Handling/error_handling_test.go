
package main

import (
	"errors"
	"strconv"
	"testing"
)

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		checkErr func(t *testing.T, err error)
	}{
		{
			name:    "valid number",
			input:   "25",
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "abc",
			wantErr: true,
			checkErr: func(t *testing.T, err error) {
				if !errors.Is(err, strconv.ErrSyntax) {
					t.Errorf("expected syntax error in chain, got: %v", err)
				}
			},
		},
		{
			name:    "negative age",
			input:   "-5",
			wantErr: true,
			checkErr: func(t *testing.T, err error) {
				var dbErr *DatabaseError
				if !errors.As(err, &dbErr) {
					t.Fatalf("expected *DatabaseError, got type: %T", err)
				}
				if dbErr.Code != 400 {
					t.Errorf("expected code 400, got %d", dbErr.Code)
				}
				if dbErr.Operation != "validate" {
					t.Errorf("expected operation 'validate', got %q", dbErr.Operation)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInput(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err != nil && tt.checkErr != nil {
				tt.checkErr(t, err)
			}
		})
	}
}
