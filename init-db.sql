CREATE TABLE IF NOT EXISTS workout (
    id SERIAL PRIMARY KEY,
    date TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS working_sets (
    id SERIAL PRIMARY KEY,
    exercise VARCHAR(255) NOT NULL,
    resistance_kg INT,
    repetitions INT,
    negative_repetitions INT,
    static_hold_seconds INT,
    workout_id INT REFERENCES workout(id)
);
