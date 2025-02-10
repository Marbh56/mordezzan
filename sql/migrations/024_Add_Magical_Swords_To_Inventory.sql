-- +goose Up
-- Add magical weapon support to character inventory
ALTER TABLE character_inventory
ADD COLUMN magical_weapon_id INTEGER REFERENCES magical_weapons (id);

-- Update the inventory view query to include magical weapon information
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
DROP VIEW IF EXISTS character_inventory_with_magic;

ALTER TABLE character_inventory
DROP COLUMN magical_weapon_id;
