-- region: SERVER

CREATE TABLE servers
(
    id               INTEGER GENERATED ALWAYS AS IDENTITY UNIQUE,
    version_major    INTEGER,
    version_minor    INTEGER,
    version_patch    INTEGER,
    database_version INTEGER
);

-- endregion

-- region: AUTH

CREATE TABLE users
(
    id              INTEGER GENERATED ALWAYS AS IDENTITY UNIQUE,
    username        VARCHAR(255) UNIQUE NOT NULL,
    name            VARCHAR(255)        NOT NULL,
    profile_picture VARCHAR(255),
    role            VARCHAR(63),
    creation_date   TIMESTAMP
);

CREATE TABLE sessions
(
    id      INTEGER GENERATED ALWAYS AS IDENTITY UNIQUE,
    user_id INTEGER             NOT NULL,
    token   VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE auth_github
(
    username VARCHAR(255) UNIQUE PRIMARY KEY,
    user_id  INTEGER NOT NULL
);

-- endregion

-- region: STORAGE

CREATE TABLE buckets
(
    id       INTEGER GENERATED ALWAYS AS IDENTITY UNIQUE,
    name     VARCHAR(255) NOT NULL,
    type     VARCHAR(63)  NOT NULL,
    max_size INTEGER
);

CREATE TABLE buckets_to_users
(
    id          INTEGER GENERATED ALWAYS AS IDENTITY UNIQUE,
    bucket_id   INTEGER     NOT NULL,
    user_id     INTEGER     NOT NULL,
    access_type VARCHAR(63) NOT NULL
);

CREATE TABLE nodes
(
    uuid        VARCHAR(63) UNIQUE PRIMARY KEY,
    parent_uuid VARCHAR(63),
    bucket_id   INTEGER      NOT NULL,
    name        VARCHAR(255) NOT NULL,
    type        VARCHAR(63)  NOT NULL,
    mime        VARCHAR(255),
    size        INTEGER
);

CREATE TABLE nodes_to_users
(
    id                     INTEGER GENERATED ALWAYS AS IDENTITY UNIQUE,
    user_id                INTEGER     NOT NULL,
    node_uuid              VARCHAR(63) NOT NULL,
    last_view_timestamp    TIMESTAMP,
    last_edition_timestamp TIMESTAMP
);

-- endregion

-- region: Foreign keys

ALTER TABLE sessions
    ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE auth_github
    ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE nodes
    ADD FOREIGN KEY (parent_uuid) REFERENCES nodes (uuid);
ALTER TABLE nodes
    ADD FOREIGN KEY (bucket_id) REFERENCES buckets (id);
ALTER TABLE buckets_to_users
    ADD FOREIGN KEY (bucket_id) REFERENCES buckets (id);
ALTER TABLE buckets_to_users
    ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE nodes_to_users
    ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE nodes_to_users
    ADD FOREIGN KEY (node_uuid) REFERENCES nodes (uuid);

-- endregion
