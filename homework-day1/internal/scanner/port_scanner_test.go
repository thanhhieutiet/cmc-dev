package scanner

import (
	"testing"

	"homework-day1/internal/model"
)

func TestPortScanner_Scan_Unauthorized(t *testing.T) {
	scanner := NewPortScanner()
	asset := &model.Asset{
		Name: "8.8.8.8",
		Type: "ip",
	}

	_, err := scanner.Scan(asset)
	if err != model.ErrPortScanUnauthorized {
		t.Errorf("expected ErrPortScanUnauthorized, got %v", err)
	}
}

func TestPortScanner_Scan_Localhost(t *testing.T) {
	scanner := NewPortScanner()
	asset := &model.Asset{
		Name: "127.0.0.1",
		Type: "ip",
	}

	result, err := scanner.Scan(asset)
	if err != nil {
		t.Fatalf("expected no error for localhost port scan, got %v", err)
	}

	if result.IPAddress != "127.0.0.1" {
		t.Errorf("expected IP 127.0.0.1, got %s", result.IPAddress)
	}
	
	if result.TotalScanned != 15 {
		t.Errorf("expected 15 ports scanned, got %d", result.TotalScanned)
	}
}

func TestPortScanner_Scan_Validation(t *testing.T) {
	scanner := NewPortScanner()
	asset := &model.Asset{
		Name: "example.com",
		Type: "domain",
	}

	_, err := scanner.Scan(asset)
	if err != model.ErrScanNotSupported {
		t.Errorf("expected ErrScanNotSupported, got %v", err)
	}
}
