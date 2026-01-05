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
		wantError bool
		wantAddr  string
	}{
		{
			name:      "valid config with custom addr",
			baseURL:   "http://localhost:8000",
			addr:      ":8080",
			wantError: false,
			wantAddr:  ":8080",
		},
		{
			name:      "valid config with default addr",
			baseURL:   "http://localhost:8000",
			addr:      "",
			wantError: false,
			wantAddr:  ":10800",
		},
		{
			name:      "missing base url",
			baseURL:   "",
			addr:      ":8080",
			wantError: true,
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

			cfg, err := NewConfig()

			if tt.wantError {
				if err == nil {
					t.Errorf("NewConfig() should return error but didn't")
				}
				return
			}

			if err != nil {
				t.Errorf("NewConfig() returned unexpected error: %v", err)
				return
			}

			if cfg.BaseUrl != tt.baseURL {
				t.Errorf("BaseUrl = %v, want %v", cfg.BaseUrl, tt.baseURL)
			}
			if cfg.Addr != tt.wantAddr {
				t.Errorf("Addr = %v, want %v", cfg.Addr, tt.wantAddr)
			}
		})
	}
}
