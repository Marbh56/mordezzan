-- +goose Up
-- Create a table to store currency conversion rates
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
DROP TABLE IF EXISTS currency_conversion_rates;