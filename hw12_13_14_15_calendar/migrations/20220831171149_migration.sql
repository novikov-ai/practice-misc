-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
                        id UUID DEFAULT uuid_generate_v4 (),
                        title VARCHAR(50),
                        description VARCHAR,
                        user_id VARCHAR,
                        date DATE,
                        duration INTERVAL,
                        notified_before INTERVAl,
                        PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS users (
                       id UUID DEFAULT uuid_generate_v4 (),
                       PRIMARY KEY (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
