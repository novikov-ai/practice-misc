START TRANSACTION;

DROP TABLE IF EXISTS events;
CREATE TABLE events (
                        id UUID DEFAULT uuid_generate_v4 (),
                        title VARCHAR(50),
                        description VARCHAR,
                        user_id VARCHAR,
                        date DATE,
                        duration INTERVAL,
                        notified_before INTERVAl,
                        PRIMARY KEY (id)
);

DROP TABLE IF EXISTS notifications;
CREATE TABLE notifications (
                               id UUID DEFAULT uuid_generate_v4 (),
                               title VARCHAR(50),
                               user_id VARCHAR,
                               date DATE,
                               PRIMARY KEY (id)
);

DROP TABLE IF EXISTS users;
CREATE TABLE users (
                       id UUID DEFAULT uuid_generate_v4 (),
                       PRIMARY KEY (id)
);

COMMIT;