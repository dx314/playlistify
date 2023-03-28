package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var (
	dbhost     = os.Getenv("PLAYLISTIFY_DBHOST")
	dbport     = os.Getenv("PLAYLISTIFY_DBPORT")
	dbuser     = os.Getenv("PLAYLISTIFY_DBUSERNAME")
	dbpassword = os.Getenv("PLAYLISTIFY_DBPASSWORD")
	dbname     = os.Getenv("PLAYLISTIFY_DBNAME")
)

// getConnection function returns a database connection
func getConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbhost,
		dbport,
		dbuser,
		dbpassword,
		dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func upsertUser(db *sql.DB, user *User) error {
	upsertQuery := `
		INSERT INTO users (spotify_id, display_name)
		VALUES ($1, $2)
		ON CONFLICT (spotify_id)
		DO UPDATE SET display_name = excluded.display_name
		RETURNING id, spotify_id, display_name, created_at;
	`

	err := db.QueryRow(upsertQuery, user.SpotifyID, user.DisplayName).Scan(&user.ID, &user.SpotifyID, &user.DisplayName, &user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func upsertToken(db *sql.DB, token *Token) error {
	upsertQuery := `
		INSERT INTO tokens (access_token, refresh_token, expires_at, user_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET access_token = excluded.access_token, refresh_token = excluded.refresh_token, expires_at = excluded.expires_at
		RETURNING access_token, refresh_token, expires_at, user_id, created_at;
	`

	err := db.QueryRow(upsertQuery, token.AccessToken, token.RefreshToken, token.ExpiresAt, token.UserID).Scan(&token.AccessToken, &token.RefreshToken, &token.ExpiresAt, &token.UserID, &token.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}
