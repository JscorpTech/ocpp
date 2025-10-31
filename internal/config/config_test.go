package config

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		addr      string
		wantPanic bool
		wantAddr  string
	}{
		{
			name:      "valid config with custom addr",
			baseURL:   "http://localhost:8000",
			addr:      ":8080",
			wantPanic: false,
			wantAddr:  ":8080",
		},
		{
			name:      "valid config with default addr",
			baseURL:   "http://localhost:8000",
			addr:      "",
			wantPanic: false,
			wantAddr:  ":10800",
		},
		{
			name:      "missing base url",
			baseURL:   "",
			addr:      ":8080",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			os.Setenv("BASE_URL", tt.baseURL)
			os.Setenv("ADDR", tt.addr)
			defer func() {
				os.Unsetenv("BASE_URL")
				os.Unsetenv("ADDR")
			}()

			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewConfig() should panic but didn't")
					}
				}()
			}

			cfg := NewConfig()

			if !tt.wantPanic {
				if cfg.BaseUrl != tt.baseURL {
					t.Errorf("BaseUrl = %v, want %v", cfg.BaseUrl, tt.baseURL)
				}
				if cfg.Addr != tt.wantAddr {
					t.Errorf("Addr = %v, want %v", cfg.Addr, tt.wantAddr)
				}
			}
		})
	}
}
