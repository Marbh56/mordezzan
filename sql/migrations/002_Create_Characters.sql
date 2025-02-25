-- sql/migrations/002_Create_Characters.sql
-- +goose Up
CREATE TABLE characters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    class TEXT NOT NULL DEFAULT 'Fighter',
    level INTEGER NOT NULL DEFAULT 1,
    max_hp INTEGER NOT NULL,
    current_hp INTEGER NOT NULL,
    strength INTEGER NOT NULL,
    dexterity INTEGER NOT NULL,
    constitution INTEGER NOT NULL,
    intelligence INTEGER NOT NULL,
    wisdom INTEGER NOT NULL,
    charisma INTEGER NOT NULL,
    experience_points INTEGER NOT NULL DEFAULT 0,
    platinum_pieces INTEGER NOT NULL DEFAULT 0,
    gold_pieces INTEGER NOT NULL DEFAULT 0,
    electrum_pieces INTEGER NOT NULL DEFAULT 0,
    silver_pieces INTEGER NOT NULL DEFAULT 0,
    copper_pieces INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE IF EXISTS characters;

-- Create a table to store currency conversion rates (from migration 018)
-- This can remain as a separate migration file but with a new number
-- Currency conversion rates table
CREATE TABLE currency_conversion_rates (
    from_currency TEXT NOT NULL,
    to_currency TEXT NOT NULL,
    rate INTEGER NOT NULL,
    PRIMARY KEY (from_currency, to_currency)
);

-- Insert the conversion rates based on Table 73
INSERT INTO
    currency_conversion_rates (from_currency, to_currency, rate)
VALUES
    -- Platinum piece conversions
    ('pp', 'pp', 1),
    ('pp', 'gp', 5),
    ('pp', 'ep', 10),
    ('pp', 'sp', 50),
    ('pp', 'cp', 250),
    -- Gold piece conversions
    ('gp', 'pp', 0), -- Use 0 for fractional conversions that need special handling
    ('gp', 'gp', 1),
    ('gp', 'ep', 2),
    ('gp', 'sp', 10),
    ('gp', 'cp', 50),
    -- Electrum piece conversions
    ('ep', 'pp', 0),
    ('ep', 'gp', 0),
    ('ep', 'ep', 1),
    ('ep', 'sp', 5),
    ('ep', 'cp', 25),
    -- Silver piece conversions
    ('sp', 'pp', 0),
    ('sp', 'gp', 0),
    ('sp', 'ep', 0),
    ('sp', 'sp', 1),
    ('sp', 'cp', 5),
    -- Copper piece conversions
    ('cp', 'pp', 0),
    ('cp', 'gp', 0),
    ('cp', 'ep', 0),
    ('cp', 'sp', 0),
    ('cp', 'cp', 1);

-- +goose Down
DROP TABLE IF EXISTS characters;
