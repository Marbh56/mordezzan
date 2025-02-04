-- +goose Up
CREATE TABLE character_weapon_masteries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    weapon_id INTEGER NOT NULL,
    mastery_level TEXT NOT NULL CHECK (mastery_level IN ('mastered', 'grand_mastery')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE,
    FOREIGN KEY (weapon_id) REFERENCES weapons (id),
    -- A character can only have one mastery level per weapon
    UNIQUE (character_id, weapon_id)
);

-- +goose Down
DROP TABLE IF EXISTS character_weapon_masteries;
