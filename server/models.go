package main

import "time"

type User struct {
	ID          string    `db:"id" json:"id"`
	SpotifyID   string    `db:"spotify_id" json:"spotify_id"`
	DisplayName string    `db:"display_name" json:"display_name"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type Token struct {
	AccessToken  string    `db:"access_token" json:"access_token"`
	RefreshToken string    `db:"refresh_token" json:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at" json:"expires_at"`
	UserID       string    `db:"user_id" json:"user_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Playlist struct {
	ID          string    `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	UserID      int       `db:"user_id" json:"user_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type Song struct {
	ID        string    `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Artist    string    `db:"artist" json:"artist"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
