package domain

import (
	"testing"
)

func TestScanType_IsValid(t *testing.T) {
	tests := []struct {
		scanType ScanType
		want     bool
	}{
		{ScanTypeDNS, true},
		{ScanTypeWHOIS, true},
		{ScanTypeSubdomain, true},
		{ScanTypeIP, true},
		{ScanTypePort, true},
		{ScanTypeSSL, true},
		{ScanTypeTech, true},
		{ScanTypeAll, true},
		{ScanType("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.scanType), func(t *testing.T) {
			if got := IsValidScanType(tt.scanType); got != tt.want {
				t.Errorf("IsValidScanType() = %v, want %v", got, tt.want)
			}
		})
	}
}
