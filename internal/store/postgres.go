// Package store reads and writes MES data in PostgreSQL.
package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/willkommlee-lgtm/mini-mes-go/internal/domain"
)

// ErrNotFound is returned when a requested record does not exist.
var ErrNotFound = errors.New("not found")

// Postgres is the PostgreSQL-backed store.
type Postgres struct {
	db *pgxpool.Pool
}

// NewPostgres creates a store that uses the given connection pool.
func NewPostgres(db *pgxpool.Pool) *Postgres {
	return &Postgres{db: db}
}

// FindMaterial looks up a material by its barcode.
func (p *Postgres) FindMaterial(ctx context.Context, barcode string) (domain.Material, error) {
	var m domain.Material
	err := p.db.QueryRow(ctx,
		"SELECT barcode, code, name FROM materials WHERE barcode = $1",
		barcode,
	).Scan(&m.Barcode, &m.Code, &m.Name)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Material{}, ErrNotFound
	}
	if err != nil {
		return domain.Material{}, err
	}
	return m, nil
}

// FindRecipeForLot returns the recipe of the work order that a lot belongs to.
func (p *Postgres) FindRecipeForLot(ctx context.Context, lotID int64) (domain.Recipe, error) {
	rows, err := p.db.Query(ctx, `
		SELECT r.id, r.name, rm.material_code
		FROM lots l
		JOIN work_orders wo      ON wo.id = l.work_order_id
		JOIN recipes r           ON r.id = wo.recipe_id
		JOIN recipe_materials rm ON rm.recipe_id = r.id
		WHERE l.id = $1`,
		lotID,
	)
	if err != nil {
		return domain.Recipe{}, err
	}
	defer rows.Close()

	var recipe domain.Recipe
	for rows.Next() {
		var code string
		if err := rows.Scan(&recipe.ID, &recipe.Name, &code); err != nil {
			return domain.Recipe{}, err
		}

		recipe.MaterialCodes = append(recipe.MaterialCodes, code)
	}
	if err := rows.Err(); err != nil {
		return domain.Recipe{}, err
	}

	if len(recipe.MaterialCodes) == 0 {
		return domain.Recipe{}, ErrNotFound
	}
	return recipe, nil
}

// SaveScan records a scan and the decision that was made.
func (p *Postgres) SaveScan(ctx context.Context, e domain.ScanEvent) error {
	_, err := p.db.Exec(ctx,
		"INSERT INTO scan_events (lot_id, barcode, result, reason) VALUES ($1, $2, $3, $4)",
		e.LotID, e.Barcode, string(e.Result), e.Reason,
	)
	return err
}
