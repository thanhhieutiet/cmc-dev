package scanner

import (
	"testing"

	"homework-day1/internal/model"
)

func TestDNSScanner_Scan_Validation(t *testing.T) {
	scanner := NewDNSScanner()
	asset := &model.Asset{
		Name: "192.168.1.1",
		Type: "ip",
	}

	_, err := scanner.Scan(asset)
	if err != model.ErrScanNotSupported {
		t.Errorf("expected ErrScanNotSupported, got %v", err)
	}
}

func TestDNSScanner_Scan_Localhost(t *testing.T) {
	scanner := NewDNSScanner()
	asset := &model.Asset{
		Name: "localhost",
		Type: "domain",
	}

	_, err := scanner.Scan(asset)
	if err != nil {
		t.Logf("DNS scan warning on localhost (expected in some offline test environments): %v", err)
	}
}
