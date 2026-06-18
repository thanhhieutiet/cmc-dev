package scanner

import (
	"testing"

	"homework-day1/internal/domain"
)

func TestSSLScanner_Scan_Validation(t *testing.T) {
	scanner := NewSSLScanner()
	asset := &domain.Asset{
		Name: "127.0.0.1",
		Type: "ip",
	}

	_, err := scanner.Scan(asset)
	if err != domain.ErrScanNotSupported {
		t.Errorf("expected ErrScanNotSupported, got %v", err)
	}
}

func TestSSLScanner_Scan_ConnectionError(t *testing.T) {
	scanner := NewSSLScanner()
	asset := &domain.Asset{
		Name: "invalid-domain-xyz123.com",
		Type: "domain",
	}

	_, err := scanner.Scan(asset)
	if err == nil {
		t.Errorf("expected connection error for invalid domain, got nil")
	}
}
