package scanner

import (
	"testing"

	"homework-day1/internal/model"
)

func TestTechScanner_Scan_Validation(t *testing.T) {
	scanner := NewTechScanner()
	asset := &model.Asset{
		Name: "127.0.0.1",
		Type: "ip",
	}

	_, err := scanner.Scan(asset)
	if err != model.ErrScanNotSupported {
		t.Errorf("expected ErrScanNotSupported, got %v", err)
	}
}

func TestTechScanner_Scan_ConnectionError(t *testing.T) {
	scanner := NewTechScanner()
	asset := &model.Asset{
		Name: "invalid-domain-xyz123.com",
		Type: "domain",
	}

	_, err := scanner.Scan(asset)
	if err == nil {
		t.Errorf("expected connection error for invalid domain, got nil")
	}
}

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		header   string
		tech     string
		expected string
	}{
		{"nginx/1.18.0", "nginx", "1.18.0"},
		{"Apache-Coyote/1.1", "Apache", "1.1"},
		{"PHP/8.0.3", "PHP", "8.0.3"},
		{"WordPress 5.7", "WordPress", "5.7"},
		{"nginx", "nginx", ""},
	}

	for _, tt := range tests {
		t.Run(tt.header+"-"+tt.tech, func(t *testing.T) {
			got := extractVersion(tt.header, tt.tech)
			if got != tt.expected {
				t.Errorf("extractVersion() = %q, want %q", got, tt.expected)
			}
		})
	}
}
