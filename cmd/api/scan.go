package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	FindRecipeForLot(ctx context.Context, lotID int64) (domain.Recipe, error)
	SaveScan(ctx context.Context, e domain.ScanEvent) error
}

// scanHandler holds what the POST /scans handler needs.
type scanHandler struct {
	materials MaterialFinder
	logger    *slog.Logger
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

	// 3. Find the recipe of the lot
	recipe, err := h.materials.FindRecipeForLot(r.Context(), req.LotID)
	if errors.Is(err, store.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("unknown lot: %d", req.LotID)})
		return
	}
	if err != nil {
		h.logger.Error("find recipe failed", "lot_id", req.LotID, "err", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	// 4. Validate against the recipe
	result, reason := domain.ValidateScan(recipe, material)

	// 5. Record the decision before answering
	event := domain.ScanEvent{
		LotID:   req.LotID,
		Barcode: req.Barcode,
		Result:  result,
		Reason:  reason,
	}
	if err := h.materials.SaveScan(r.Context(), event); err != nil {
		h.logger.Error("save scan failed", "lot_id", req.LotID, "barcode", req.Barcode, "err", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	// 6. Answer the PC
	writeJSON(w, http.StatusOK, scanResponse{Result: result, Reason: reason})
}
