-- +goose Up
-- Migration to add magic items to the database

-- Add magical weapon variants
INSERT INTO weapons (
    name, 
    reach, 
    cost_gp, 
    weight, 
    range_short, 
    range_medium, 
    range_long, 
    attacks_per_round, 
    damage, 
    enhancement_bonus
)
SELECT 
    name || ' +1', 
    reach, 
    cost_gp + 2000, 
    weight, 
    range_short, 
    range_medium, 
    range_long, 
    attacks_per_round, 
    damage, 
    1
FROM 
    weapons
WHERE 
    enhancement_bonus IS NULL OR enhancement_bonus = 0;

INSERT INTO weapons (
    name, 
    reach, 
    cost_gp, 
    weight, 
    range_short, 
    range_medium, 
    range_long, 
    attacks_per_round, 
    damage, 
    enhancement_bonus
)
SELECT 
    name || ' +2', 
    reach, 
    cost_gp + 8000, 
    weight, 
    range_short, 
    range_medium, 
    range_long, 
    attacks_per_round, 
    damage, 
    2
FROM 
    weapons
WHERE 
    enhancement_bonus IS NULL OR enhancement_bonus = 0;

INSERT INTO weapons (
    name, 
    reach, 
    cost_gp, 
    weight, 
    range_short, 
    range_medium, 
    range_long, 
    attacks_per_round, 
    damage, 
    enhancement_bonus
)
SELECT 
    name || ' +3', 
    reach, 
    cost_gp + 18000, 
    weight, 
    range_short, 
    range_medium, 
    range_long, 
    attacks_per_round, 
    damage, 
    3
FROM 
    weapons
WHERE 
    enhancement_bonus IS NULL OR enhancement_bonus = 0;

-- Add magical armor variants
INSERT INTO armor (
    name, 
    armor_class, 
    cost_gp, 
    damage_reduction, 
    weight, 
    armor_type, 
    movement_rate, 
    enhancement_bonus
)
SELECT 
    name || ' +1', 
    armor_class - 1, -- Improve AC by 1
    cost_gp + 2000, 
    damage_reduction, 
    weight, 
    armor_type, 
    movement_rate, 
    1
FROM 
    armor
WHERE 
    enhancement_bonus IS NULL OR enhancement_bonus = 0;

INSERT INTO armor (
    name, 
    armor_class, 
    cost_gp, 
    damage_reduction, 
    weight, 
    armor_type, 
    movement_rate, 
    enhancement_bonus
)
SELECT 
    name || ' +2', 
    armor_class - 2, -- Improve AC by 2
    cost_gp + 8000, 
    damage_reduction, 
    weight, 
    armor_type, 
    movement_rate, 
    2
FROM 
    armor
WHERE 
    enhancement_bonus IS NULL OR enhancement_bonus = 0;

INSERT INTO armor (
    name, 
    armor_class, 
    cost_gp, 
    damage_reduction, 
    weight, 
    armor_type, 
    movement_rate, 
    enhancement_bonus
)
SELECT 
    name || ' +3', 
    armor_class - 3, -- Improve AC by 3
    cost_gp + 18000, 
    damage_reduction, 
    weight, 
    armor_type, 
    movement_rate, 
    3
FROM 
    armor
WHERE 
    enhancement_bonus IS NULL OR enhancement_bonus = 0;

-- +goose Down
-- Remove magical items

-- Delete +3 weapons
DELETE FROM weapons WHERE name LIKE '% +3';

-- Delete +2 weapons
DELETE FROM weapons WHERE name LIKE '% +2';

-- Delete +1 weapons
DELETE FROM weapons WHERE name LIKE '% +1';

-- Delete +3 armor
DELETE FROM armor WHERE name LIKE '% +3';

-- Delete +2 armor
DELETE FROM armor WHERE name LIKE '% +2';

-- Delete +1 armor
DELETE FROM armor WHERE name LIKE '% +1';