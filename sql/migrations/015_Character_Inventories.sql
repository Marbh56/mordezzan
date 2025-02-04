-- +goose Up
-- Create a polymorphic inventory items table that can reference different types of equipment
CREATE TABLE character_inventory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    -- The type of item (equipment, weapon, armor, ammunition, container)
    item_type TEXT NOT NULL CHECK (
        item_type IN (
            'equipment',
            'weapon',
            'armor',
            'ammunition',
            'container',
            'shield',
            'ranged_weapon'
        )
    ),
    -- The ID in the corresponding table
    item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    -- If the item is stored in a container, reference that container's inventory entry
    container_inventory_id INTEGER,
    -- If the item is equipped, reference which slot it's in
    equipment_slot_id INTEGER,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE,
    FOREIGN KEY (container_inventory_id) REFERENCES character_inventory (id),
    FOREIGN KEY (equipment_slot_id) REFERENCES equipment_slots (id),
    -- Composite index for efficient character inventory queries
    UNIQUE (
        character_id,
        item_type,
        item_id,
        container_inventory_id
    )
);

-- Create indexes for better query performance
CREATE INDEX idx_character_inventory_character_id ON character_inventory (character_id);

CREATE INDEX idx_character_inventory_container ON character_inventory (container_inventory_id);

CREATE INDEX idx_character_inventory_equipment_slot ON character_inventory (equipment_slot_id);

-- +goose Down
DROP TABLE IF EXISTS character_inventory;
