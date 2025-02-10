-- +goose Up
-- First drop the view that depends on the table
DROP VIEW IF EXISTS character_inventory_with_magic;

-- Now we can modify the table
CREATE TABLE character_inventory_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    item_type TEXT NOT NULL CHECK (
        item_type IN (
            'equipment',
            'weapon',
            'armor',
            'ammunition',
            'container',
            'shield',
            'ranged_weapon',
            'magical_weapon'
        )
    ),
    item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    container_inventory_id INTEGER,
    equipment_slot_id INTEGER,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    magical_weapon_id INTEGER REFERENCES magical_weapons (id),
    FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE,
    FOREIGN KEY (container_inventory_id) REFERENCES character_inventory (id),
    FOREIGN KEY (equipment_slot_id) REFERENCES equipment_slots (id)
);

INSERT INTO
    character_inventory_new
SELECT
    *
FROM
    character_inventory;

DROP TABLE character_inventory;

ALTER TABLE character_inventory_new
RENAME TO character_inventory;

-- Recreate the indexes
CREATE INDEX idx_character_inventory_character_id ON character_inventory (character_id);

CREATE INDEX idx_character_inventory_container ON character_inventory (container_inventory_id);

CREATE INDEX idx_character_inventory_equipment_slot ON character_inventory (equipment_slot_id);

-- Recreate the view
CREATE VIEW character_inventory_with_magic AS
SELECT
    ci.*,
    CASE
        WHEN ci.magical_weapon_id IS NOT NULL THEN w.name || ' +' || mw.enhancement_bonus
        WHEN ci.item_type = 'weapon' THEN w.name
        WHEN ci.item_type = 'armor' THEN a.name
        WHEN ci.item_type = 'ammunition' THEN am.name
        WHEN ci.item_type = 'container' THEN c.name
        WHEN ci.item_type = 'shield' THEN s.name
        WHEN ci.item_type = 'ranged_weapon' THEN rw.name
        WHEN ci.item_type = 'equipment' THEN e.name
    END as item_name,
    CASE
        WHEN ci.item_type = 'weapon' THEN w.weight
        WHEN ci.item_type = 'armor' THEN a.weight
        WHEN ci.item_type = 'ammunition' THEN am.weight
        WHEN ci.item_type = 'container' THEN c.weight
        WHEN ci.item_type = 'shield' THEN s.weight
        WHEN ci.item_type = 'ranged_weapon' THEN rw.weight
        WHEN ci.item_type = 'equipment' THEN e.weight
        ELSE 0
    END as item_weight
FROM
    character_inventory ci
    LEFT JOIN weapons w ON ci.item_type = 'weapon'
    AND ci.item_id = w.id
    LEFT JOIN magical_weapons mw ON ci.magical_weapon_id = mw.id
    LEFT JOIN armor a ON ci.item_type = 'armor'
    AND ci.item_id = a.id
    LEFT JOIN ammunition am ON ci.item_type = 'ammunition'
    AND ci.item_id = am.id
    LEFT JOIN containers c ON ci.item_type = 'container'
    AND ci.item_id = c.id
    LEFT JOIN shields s ON ci.item_type = 'shield'
    AND ci.item_id = s.id
    LEFT JOIN ranged_weapons rw ON ci.item_type = 'ranged_weapon'
    AND ci.item_id = rw.id
    LEFT JOIN equipment e ON ci.item_type = 'equipment'
    AND ci.item_id = e.id;

-- +goose Down
-- First drop the view
DROP VIEW IF EXISTS character_inventory_with_magic;

-- Then modify the table
CREATE TABLE character_inventory_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
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
    item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    container_inventory_id INTEGER,
    equipment_slot_id INTEGER,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    magical_weapon_id INTEGER REFERENCES magical_weapons (id),
    FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE,
    FOREIGN KEY (container_inventory_id) REFERENCES character_inventory (id),
    FOREIGN KEY (equipment_slot_id) REFERENCES equipment_slots (id)
);

INSERT INTO
    character_inventory_new
SELECT
    *
FROM
    character_inventory
WHERE
    item_type != 'magical_weapon';

DROP TABLE character_inventory;

ALTER TABLE character_inventory_new
RENAME TO character_inventory;

-- Recreate the indexes
CREATE INDEX idx_character_inventory_character_id ON character_inventory (character_id);

CREATE INDEX idx_character_inventory_container ON character_inventory (container_inventory_id);

CREATE INDEX idx_character_inventory_equipment_slot ON character_inventory (equipment_slot_id);

-- Recreate the view
CREATE VIEW character_inventory_with_magic AS
SELECT
    ci.*,
    CASE
        WHEN ci.magical_weapon_id IS NOT NULL THEN w.name || ' +' || mw.enhancement_bonus
        WHEN ci.item_type = 'weapon' THEN w.name
        WHEN ci.item_type = 'armor' THEN a.name
        WHEN ci.item_type = 'ammunition' THEN am.name
        WHEN ci.item_type = 'container' THEN c.name
        WHEN ci.item_type = 'shield' THEN s.name
        WHEN ci.item_type = 'ranged_weapon' THEN rw.name
        WHEN ci.item_type = 'equipment' THEN e.name
    END as item_name,
    CASE
        WHEN ci.item_type = 'weapon' THEN w.weight
        WHEN ci.item_type = 'armor' THEN a.weight
        WHEN ci.item_type = 'ammunition' THEN am.weight
        WHEN ci.item_type = 'container' THEN c.weight
        WHEN ci.item_type = 'shield' THEN s.weight
        WHEN ci.item_type = 'ranged_weapon' THEN rw.weight
        WHEN ci.item_type = 'equipment' THEN e.weight
        ELSE 0
    END as item_weight
FROM
    character_inventory ci
    LEFT JOIN weapons w ON ci.item_type = 'weapon'
    AND ci.item_id = w.id
    LEFT JOIN magical_weapons mw ON ci.magical_weapon_id = mw.id
    LEFT JOIN armor a ON ci.item_type = 'armor'
    AND ci.item_id = a.id
    LEFT JOIN ammunition am ON ci.item_type = 'ammunition'
    AND ci.item_id = am.id
    LEFT JOIN containers c ON ci.item_type = 'container'
    AND ci.item_id = c.id
    LEFT JOIN shields s ON ci.item_type = 'shield'
    AND ci.item_id = s.id
    LEFT JOIN ranged_weapons rw ON ci.item_type = 'ranged_weapon'
    AND ci.item_id = rw.id
    LEFT JOIN equipment e ON ci.item_type = 'equipment'
    AND ci.item_id = e.id;
