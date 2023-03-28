CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id           uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    spotify_id   varchar(255) UNIQUE NOT NULL,
    display_name varchar(255),
    email        varchar(255) UNIQUE NOT NULL,
    created_at   timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tokens
(
    id            uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    access_token  varchar(255)             NOT NULL,
    token_type    varchar(255)             NOT NULL,
    expires_at    timestamp NOT NULL,
    refresh_token varchar(255)             NOT NULL,
    user_id       uuid REFERENCES users (id),
    created_at    timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE playlists
(
    id          uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    title       varchar(255) NOT NULL,
    description text,
    user_id     uuid REFERENCES users (id),
    created_at  timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE songs
(
    id         uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    title      varchar(255) NOT NULL,
    artist     varchar(255) NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE playlist_songs
(
    id          uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    playlist_id uuid REFERENCES playlists (id),
    song_id     uuid REFERENCES songs (id),
    created_at  timestamp DEFAULT CURRENT_TIMESTAMP
);