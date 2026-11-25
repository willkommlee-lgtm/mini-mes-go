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
