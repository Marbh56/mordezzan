-- +goose Up
CREATE TABLE coins (
    denomination TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    weight_per_coin DECIMAL(10, 3) NOT NULL, -- Weight in pounds
    base_value INTEGER NOT NULL -- Value in copper pieces
);

INSERT INTO
    coins (denomination, name, weight_per_coin, base_value)
VALUES
    ('pp', 'Platinum Piece', 0.01, 500), -- 1 pp = 5 gp = 500 cp
    ('gp', 'Gold Piece', 0.01, 100), -- 1 gp = 100 cp
    ('ep', 'Electrum Piece', 0.01, 50), -- 1 ep = 5 sp = 50 cp
    ('sp', 'Silver Piece', 0.01, 10), -- 1 sp = 10 cp
    ('cp', 'Copper Piece', 0.01, 1);

-- Base unit, 100 coins = 1 pound
-- +goose Down
DROP TABLE IF EXISTS coins;
