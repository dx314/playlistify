-- Write your 'up' migration SQL here
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create schema if not exists public;

CREATE TABLE users
(
    id           uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    spotify_id   text UNIQUE NOT NULL,
    display_name varchar(255),
    created_at   timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tokens
(
    access_token  text             NOT NULL,
    expires_at    timestamp NOT NULL,
    refresh_token text             NOT NULL,
    user_id       uuid PRIMARY KEY REFERENCES users (id),
    created_at    timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE playlists
(
    id          uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    title       text NOT NULL,
    description text,
    user_id     uuid REFERENCES users (id),
    created_at  timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE songs
(
    id         uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    title      text NOT NULL,
    artist     text NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE playlist_songs
(
    id          uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    playlist_id uuid REFERENCES playlists (id),
    song_id     uuid REFERENCES songs (id),
    created_at  timestamp DEFAULT CURRENT_TIMESTAMP
);