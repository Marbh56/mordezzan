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

CREATE TABLE container_allowed_types (
    container_id INTEGER NOT NULL,
    item_type TEXT NOT NULL,
    ammo_type TEXT,
    FOREIGN KEY (container_id) REFERENCES containers (id) ON DELETE CASCADE,
    PRIMARY KEY (container_id, item_type)
);

-- +goose Down
DROP TABLE IF EXISTS container_allowed_types;

DROP TABLE IF EXISTS containers;
