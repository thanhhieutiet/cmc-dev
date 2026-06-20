package scanner

import (
	"testing"

	"homework-day1/internal/model"
)

func TestIPScanner_Scan_PrivateIP(t *testing.T) {
	scanner := NewIPScanner()
	asset := &model.Asset{
		Name: "127.0.0.1",
		Type: "ip",
	}

	result, err := scanner.Scan(asset)
	if err != nil {
		t.Fatalf("expected no error for private IP, got %v", err)
	}

	if result.IPAddress != "127.0.0.1" {
		t.Errorf("expected IP 127.0.0.1, got %s", result.IPAddress)
	}

	if result.Geolocation.City != "Localhost" {
		t.Errorf("expected Geolocation City to be Localhost, got %s", result.Geolocation.City)
	}

	if result.ASN.Number != 0 {
		t.Errorf("expected ASN number 0, got %d", result.ASN.Number)
	}
}

func TestIPScanner_Scan_Validation(t *testing.T) {
	scanner := NewIPScanner()
	asset := &model.Asset{
		Name: "example.com",
		Type: "domain",
	}

	_, err := scanner.Scan(asset)
	if err != model.ErrScanNotSupported {
		t.Errorf("expected ErrScanNotSupported, got %v", err)
	}
}
