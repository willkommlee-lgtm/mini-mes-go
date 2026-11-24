package main

import (
	"encoding/json"
	"net/http"

	"github.com/willkommlee-lgtm/mini-mes-go/internal/domain"
)

// scanRequest is what a shop-floor PC sends after scanning a material.
type scanRequest struct {
	LotID   int64  `json:"lot_id"`
	Barcode string `json:"barcode"`
}

// scanResponse tells the PC whether the material may enter the mixer.
type scanResponse struct {
	Result domain.ScanResult `json:"result"`
	Reason string            `json:"reason,omitempty"`
}

// Temporary in-memory data. Replaced by PostgreSQL in the next step.
var materials = map[string]domain.Material{
	"MAT-0001": {Barcode: "MAT-0001", Code: "NR-01", Name: "Natural rubber"},
	"MAT-0002": {Barcode: "MAT-0002", Code: "CB-02", Name: "Carbon black"},
	"MAT-0003": {Barcode: "MAT-0003", Code: "OIL-03", Name: "Process oil"},
	"MAT-9999": {Barcode: "MAT-9999", Code: "SI-09", Name: "Silica"},
}

var compoundA = domain.Recipe{
	Name:          "Compound-A",
	MaterialCodes: []string{"NR-01", "CB-02", "OIL-03"},
}

// handleScan handles POST /scans.
func handleScan(w http.ResponseWriter, r *http.Request) {
	// 1. Read the JSON body into req
	var req scanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	// 2. Find the material by barcode
	material, ok := materials[req.Barcode]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "unknown barcode: " + req.Barcode})
		return
	}
	// - materials 맵에서 req.Barcode로 찾기 (4강의 value, ok 패턴)
	// - 없으면: writeJSON(w, http.StatusNotFound, map[string]string{"error": "unknown barcode: " + req.Barcode})
	//          그리고 return

	// 3. Validate against the recipe
	result, reason := domain.ValidateScan(compoundA, material)
	// - domain.ValidateScan을 compoundA와 찾은 material로 호출해서 result, reason 받기

	// 4. Answer the PC
	writeJSON(w, http.StatusOK, scanResponse{Result: result, Reason: reason})
}
