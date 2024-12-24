-- +goose Up
CREATE TABLE characters (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL,
    class VARCHAR(50) NOT NULL,
    race VARCHAR(50) NOT NULL,
    alignment VARCHAR(50) NOT NULL,
    level INTEGER NOT NULL DEFAULT 1,
    hit_points_total INTEGER NOT NULL,
    hit_points_current INTEGER NOT NULL CHECK (hit_points_current >= -10),
    languages TEXT,
    religion VARCHAR(100),
    secondary_skills TEXT,
    place_of_origin VARCHAR(255),
    gender VARCHAR(50),
    age INTEGER,
    height VARCHAR(20),
    weight INTEGER,
    eye_color VARCHAR(50),
    hair_color VARCHAR(50),
    complexion VARCHAR(100),
    xp_current INTEGER DEFAULT 0,
    xp_bonus BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE characters;
