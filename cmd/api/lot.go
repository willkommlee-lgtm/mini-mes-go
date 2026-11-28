package main

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/willkommlee-lgtm/mini-mes-go/internal/domain"
)

// GenealogyReader lists the materials that went into a lot.
type GenealogyReader interface {
	LotGenealogy(ctx context.Context, lotID int64) ([]domain.GenealogyEntry, error)
}

// lotHandler holds what the lot endpoints need.
type lotHandler struct {
	genealogy GenealogyReader
	logger    *slog.Logger
}

// handleGenealogy handles GET /lots/{id}/genealogy.
func (h *lotHandler) handleGenealogy(w http.ResponseWriter, r *http.Request) {
	// 1. Read the lot id from the URL
	lotID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "lot id must be a number"})
		return
	}

	// 2. Load the genealogy
	entries, err := h.genealogy.LotGenealogy(r.Context(), lotID)
	if err != nil {
		h.logger.Error("load genealogy failed", "lot_id", lotID, "err", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	// 3. Answer
	writeJSON(w, http.StatusOK, entries)
}
