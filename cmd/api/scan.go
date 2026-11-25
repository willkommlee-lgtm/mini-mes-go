package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/willkommlee-lgtm/mini-mes-go/internal/domain"
	"github.com/willkommlee-lgtm/mini-mes-go/internal/store"
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

// MaterialFinder finds a material by barcode.
// The real one reads PostgreSQL; a test can use a fake one.
type MaterialFinder interface {
	FindMaterial(ctx context.Context, barcode string) (domain.Material, error)
}

// scanHandler holds what the POST /scans handler needs.
type scanHandler struct {
	materials MaterialFinder
	logger    *slog.Logger
}

// Temporary recipe. Moved to PostgreSQL in step 3-3.
var compoundA = domain.Recipe{
	Name:          "Compound-A",
	MaterialCodes: []string{"NR-01", "CB-02", "OIL-03"},
}

// handleScan handles POST /scans.
func (h *scanHandler) handleScan(w http.ResponseWriter, r *http.Request) {
	// 1. Read the JSON body into req
	var req scanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	// 2. Find the material by barcode (now from PostgreSQL)
	material, err := h.materials.FindMaterial(r.Context(), req.Barcode)

	if errors.Is(err, store.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "unknown barcode: " + req.Barcode})
		return
	}
	if err != nil {
		h.logger.Error("find material failed", "barcode", req.Barcode, "err", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	// - err가 store.ErrNotFound면: 404 "unknown barcode: ..." 응답하고 return
	// - 그 밖의 err면: h.logger.Error로 기록하고, 500 "internal error" 응답하고 return

	// 3. Validate against the recipe
	result, reason := domain.ValidateScan(compoundA, material)

	// 4. Answer the PC
	writeJSON(w, http.StatusOK, scanResponse{Result: result, Reason: reason})
}
