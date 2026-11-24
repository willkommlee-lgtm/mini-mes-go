package domain

import "testing"

func TestValidateScan(t *testing.T) {
	recipe := Recipe{
		Name:          "Compound-A",
		MaterialCodes: []string{"NR-01", "CB-02", "OIL-03"},
	}

	tests := []struct {
		name       string
		material   Material
		wantResult ScanResult
	}{
		{"first material in recipe", Material{Barcode: "MAT-0001", Code: "NR-01"}, ScanAllow},
		{"middle material in recipe", Material{Barcode: "MAT-0002", Code: "CB-02"}, ScanAllow},
		{"last material in recipe", Material{Barcode: "MAT-0003", Code: "OIL-03"}, ScanAllow},
		{"material not in recipe", Material{Barcode: "MAT-9999", Code: "SI-09"}, ScanBlock},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotReason := ValidateScan(recipe, tt.material)

			if gotResult != tt.wantResult {
				t.Errorf("ValidateScan() result = %s, want %s", gotResult, tt.wantResult)
			}
			if gotResult == ScanBlock && gotReason == "" {
				t.Errorf("ValidateScan() blocked without a reason")
			}
			if gotResult == ScanAllow && gotReason != "" {
				t.Errorf("ValidateScan() allowed but gave a reason: %q", gotReason)
			}
		})
	}
}
