package model

import (
	"testing"
)

func TestAsset_Validate(t *testing.T) {
	tests := []struct {
		name    string
		asset   Asset
		wantErr error
	}{
		{
			name: "valid asset",
			asset: Asset{
				Name:   "example.com",
				Type:   "domain",
				Status: "active",
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			asset: Asset{
				Name:   "",
				Type:   "domain",
				Status: "active",
			},
			wantErr: ErrInvalidName,
		},
		{
			name: "invalid type",
			asset: Asset{
				Name:   "example.com",
				Type:   "invalid_type",
				Status: "active",
			},
			wantErr: ErrInvalidAssetType,
		},
		{
			name: "invalid status",
			asset: Asset{
				Name:   "example.com",
				Type:   "domain",
				Status: "invalid_status",
			},
			wantErr: ErrInvalidAssetStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.asset.Validate()
			if err != tt.wantErr {
				t.Errorf("Asset.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
