-- +goose Up
CREATE TABLE containers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    cost_gp DECIMAL(10, 2) NOT NULL,
    weight INTEGER,
    capacity_weight INTEGER NOT NULL,
    capacity_items INTEGER,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS containers;
