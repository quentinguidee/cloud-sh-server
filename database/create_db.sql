-- region: AUTH

CREATE TABLE IF NOT EXISTS users
(
    id              INTEGER GENERATED ALWAYS AS IDENTITY UNIQUE,
    username        VARCHAR(255) UNIQUE,
    name            VARCHAR(255),
    profile_picture VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS sessions
(
    id      INTEGER GENERATED ALWAYS AS IDENTITY UNIQUE,
    user_id INTEGER,
    token   VARCHAR(255) UNIQUE
);

CREATE TABLE IF NOT EXISTS auth_github
(
    username VARCHAR(255) UNIQUE PRIMARY KEY,
    user_id  INTEGER
);

-- endregion

-- region: STORAGE

CREATE TABLE IF NOT EXISTS buckets
(
    id        INTEGER GENERATED ALWAYS AS IDENTITY UNIQUE,
    name      VARCHAR(255),
    root_node VARCHAR(63),
    type      VARCHAR(63)
);

CREATE TABLE IF NOT EXISTS buckets_access
(
    id          INTEGER GENERATED ALWAYS AS IDENTITY UNIQUE,
    bucket_id   INTEGER,
    user_id     INTEGER,
    access_type VARCHAR(63)
);

CREATE TABLE IF NOT EXISTS buckets_nodes
(
    uuid      VARCHAR(63) UNIQUE PRIMARY KEY,
    name      VARCHAR(255),
    type      VARCHAR(63),
    mime      VARCHAR(63),
    size      INTEGER,
    bucket_id INTEGER
);

CREATE TABLE IF NOT EXISTS buckets_nodes_associations
(
    id        INTEGER GENERATED ALWAYS AS IDENTITY UNIQUE,
    from_node VARCHAR(63),
    to_node   VARCHAR(63)
);

-- endregion

-- region: Foreign keys

ALTER TABLE sessions
    ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE auth_github
    ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE buckets
    ADD FOREIGN KEY (root_node) REFERENCES buckets_nodes (uuid);
ALTER TABLE buckets_access
    ADD FOREIGN KEY (bucket_id) REFERENCES buckets (id);
ALTER TABLE buckets_access
    ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE buckets_nodes
    ADD FOREIGN KEY (bucket_id) REFERENCES buckets (id);
ALTER TABLE buckets_nodes_associations
    ADD FOREIGN KEY (from_node) REFERENCES buckets_nodes (uuid);
ALTER TABLE buckets_nodes_associations
    ADD FOREIGN KEY (to_node) REFERENCES buckets_nodes (uuid);

-- endregion
