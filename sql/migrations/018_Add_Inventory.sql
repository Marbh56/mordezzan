-- +goose Up
-- Create new items table
CREATE TABLE items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    weight REAL NOT NULL DEFAULT 0,
    value REAL NOT NULL DEFAULT 0,
    stackable BOOLEAN NOT NULL DEFAULT 0,
    max_stack INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create item tags table
CREATE TABLE item_tags (
    item_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY (item_id, tag),
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE
);

-- Create improved containers table
CREATE TABLE containers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    base_item_id INTEGER NOT NULL,
    capacity REAL NOT NULL,
    max_items INTEGER,
    container_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (base_item_id) REFERENCES items (id) ON DELETE CASCADE
);

-- Create container allowed tags table
CREATE TABLE container_allowed_tags (
    container_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY (container_id, tag),
    FOREIGN KEY (container_id) REFERENCES containers (id) ON DELETE CASCADE
);

-- Create improved character inventory table
CREATE TABLE character_inventory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    item_id INTEGER NOT NULL,
    item_type TEXT NOT NULL CHECK (
        item_type IN (
            'equipment',
            'weapon',
            'armor',
            'ammunition',
            'container',
            'shield',
            'ranged_weapon',
            'consumable',
            'magical_item'
        )
    ),
    quantity INTEGER NOT NULL DEFAULT 1,
    container_id INTEGER,
    slot_id INTEGER,
    position TEXT,
    custom_name TEXT,
    custom_notes TEXT,
    is_identified BOOLEAN NOT NULL DEFAULT 1,
    charges INTEGER,
    condition TEXT NOT NULL DEFAULT 'good',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE,
    FOREIGN KEY (container_id) REFERENCES character_inventory (id) ON DELETE SET NULL,
    FOREIGN KEY (slot_id) REFERENCES equipment_slots (id) ON DELETE SET NULL,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_character_inventory_character_id ON character_inventory (character_id);
CREATE INDEX idx_character_inventory_container_id ON character_inventory (container_id);
CREATE INDEX idx_character_inventory_slot_id ON character_inventory (slot_id);
CREATE INDEX idx_character_inventory_item_type ON character_inventory (item_type);
CREATE INDEX idx_item_tags_tag ON item_tags (tag);

-- +goose Down
-- Drop all tables in reverse order
DROP TABLE IF EXISTS character_inventory;
DROP TABLE IF EXISTS container_allowed_tags;
DROP TABLE IF EXISTS containers;
DROP TABLE IF EXISTS item_tags;
DROP TABLE IF EXISTS items;