package scanner

import (
	"testing"

	"homework-day1/internal/domain"
)

func TestDNSScanner_Scan_Validation(t *testing.T) {
	scanner := NewDNSScanner()
	asset := &domain.Asset{
		Name: "192.168.1.1",
		Type: "ip",
	}

	_, err := scanner.Scan(asset)
	if err != domain.ErrScanNotSupported {
		t.Errorf("expected ErrScanNotSupported, got %v", err)
	}
}

func TestDNSScanner_Scan_Localhost(t *testing.T) {
	scanner := NewDNSScanner()
	asset := &domain.Asset{
		Name: "localhost",
		Type: "domain",
	}

	_, err := scanner.Scan(asset)
	if err != nil {
		t.Logf("DNS scan warning on localhost (expected in some offline test environments): %v", err)
	}
}
