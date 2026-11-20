-- mini-mes-go: initial schema (runs automatically when the Postgres container starts for the first time)

CREATE TABLE recipes (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE
);

-- Which material codes a recipe allows
CREATE TABLE recipe_materials (
    recipe_id     BIGINT NOT NULL REFERENCES recipes(id),
    material_code TEXT   NOT NULL,
    PRIMARY KEY (recipe_id, material_code)
);

CREATE TABLE materials (
    barcode     TEXT PRIMARY KEY,
    code        TEXT NOT NULL,
    name        TEXT NOT NULL
);

CREATE TABLE work_orders (
    id          BIGSERIAL PRIMARY KEY,
    product     TEXT    NOT NULL,
    recipe_id   BIGINT  NOT NULL REFERENCES recipes(id),
    quantity    INT     NOT NULL CHECK (quantity > 0),
    status      TEXT    NOT NULL DEFAULT 'PLANNED',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE lots (
    id             BIGSERIAL PRIMARY KEY,
    work_order_id  BIGINT NOT NULL REFERENCES work_orders(id),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE scan_events (
    id          BIGSERIAL PRIMARY KEY,
    lot_id      BIGINT NOT NULL REFERENCES lots(id),
    barcode     TEXT   NOT NULL,
    result      TEXT   NOT NULL CHECK (result IN ('ALLOW', 'BLOCK')),
    reason      TEXT,
    scanned_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_scan_events_lot ON scan_events(lot_id);

-- Sample data for local testing
INSERT INTO recipes (name) VALUES ('Compound-A');
INSERT INTO recipe_materials (recipe_id, material_code) VALUES (1, 'NR-01'), (1, 'CB-02'), (1, 'OIL-03');
INSERT INTO materials (barcode, code, name) VALUES
    ('MAT-0001', 'NR-01',  'Natural rubber'),
    ('MAT-0002', 'CB-02',  'Carbon black'),
    ('MAT-0003', 'OIL-03', 'Process oil'),
    ('MAT-9999', 'SI-09',  'Silica (not in Compound-A)');
