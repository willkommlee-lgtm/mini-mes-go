// Package domain holds the core MES concepts of mini-mes-go.
//
// The model is a simplified version of a rubber mixing process:
// a work order follows a recipe, each production batch is a lot,
// and every raw material is scanned before it may enter the mixer.
package domain

import "time"

// WorkOrder is a production order for one product, following one recipe.
type WorkOrder struct {
	ID        int64     `json:"id"`
	Product   string    `json:"product"`
	RecipeID  int64     `json:"recipe_id"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"` // e.g. PLANNED, RUNNING, DONE
	CreatedAt time.Time `json:"created_at"`
}

// Material is a raw material identified by its barcode.
type Material struct {
	Barcode string `json:"barcode"`
	Code    string `json:"code"` // material type, e.g. "NR-01" (natural rubber)
	Name    string `json:"name"`
}

// Lot is one production batch of a work order.
type Lot struct {
	ID          int64     `json:"id"`
	WorkOrderID int64     `json:"work_order_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// ScanResult is the decision for a scanned material.
type ScanResult string

const (
	ScanAllow ScanResult = "ALLOW"
	ScanBlock ScanResult = "BLOCK"
)

// ScanEvent records one barcode scan at the mixer and its decision.
type ScanEvent struct {
	ID        int64      `json:"id"`
	LotID     int64      `json:"lot_id"`
	Barcode   string     `json:"barcode"`
	Result    ScanResult `json:"result"`
	Reason    string     `json:"reason,omitempty"` // why it was blocked
	ScannedAt time.Time  `json:"scanned_at"`
}

// Recipe lists the material codes a compound is allowed to contain.
type Recipe struct {
	ID            int64
	Name          string
	MaterialCodes []string
}

func ValidateScan(recipe Recipe, material Material) (ScanResult, string) {
	for _, c := range recipe.MaterialCodes {
		if c == material.Code {
			return ScanAllow, ""
		}
	}
	return ScanBlock, "material " + material.Code + " is not in recipe " + recipe.Name
}

// GenealogyEntry is one material that went into a lot.
type GenealogyEntry struct {
	Barcode      string    `json:"barcode"`
	MaterialCode string    `json:"material_code"`
	MaterialName string    `json:"material_name"`
	ScannedAt    time.Time `json:"scanned_at"`
}
