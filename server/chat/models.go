package main

import "time"

type User struct {
	ID          string    `db:"id"`
	SpotifyID   string    `db:"spotify_id"`
	DisplayName string    `db:"display_name"`
	Email       string    `db:"email"`
	CreatedAt   time.Time `db:"created_at"`
}

type Token struct {
	ID           int       `db:"id"`
	AccessToken  string    `db:"access_token"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	UserID       int       `db:"user_id"`
	CreatedAt    time.Time `db:"created_at"`
}

type Playlist struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	UserID      int       `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
}

type Song struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Artist    string    `db:"artist"`
	CreatedAt time.Time `db:"created_at"`
}
