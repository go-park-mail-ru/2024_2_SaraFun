CREATE TABLE IF NOT EXISTS award (
    id SERIAL PRIMARY KEY,
    day_number INT NOT NULL UNIQUE CHECK (day_number BETWEEN 0 and 6),
    award_type text NOT NULL,
    award_count INT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_activity (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    last_login TIMESTAMP NOT NULL,
    consecutive_days INT NOT NULL DEFAULT 0,

    CONSTRAINT fk_user FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

INSERT INTO user_activity (user_id, last_login, consecutive_days)
    SELECT id, NOW(), 0 FROM users;

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE award TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE user_activity TO app_user;
GRANT USAGE, SELECT ON SEQUENCE award_id_seq TO app_user;
GRANT USAGE, SELECT ON SEQUENCE user_activity_id_seq TO app_user;