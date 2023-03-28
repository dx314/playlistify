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

	fmt.Println(connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func upsertUser(db *sql.DB, user *User) (*User, error) {
	upsertQuery := `
		INSERT INTO users (spotify_id, display_name, email)
		VALUES ($1, $2, $3)
		ON CONFLICT (spotify_id)
		DO UPDATE SET display_name = excluded.display_name, email = excluded.email
		RETURNING id, spotify_id, display_name, email;
	`

	err := db.QueryRow(upsertQuery, user.SpotifyID, user.DisplayName, user.Email).Scan(&user.ID, &user.SpotifyID, &user.DisplayName, &user.Email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func upsertToken(db *sql.DB, token *Token) (*Token, error) {
	upsertQuery := `
		INSERT INTO tokens (access_token, refresh_token, expires_at, user_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET access_token = excluded.access_token, refresh_token = excluded.refresh_token, expires_at = excluded.expires_at
		RETURNING id, access_token, refresh_token, expires_at, user_id;
	`

	err := db.QueryRow(upsertQuery, token.AccessToken, token.RefreshToken, token.ExpiresAt, token.UserID).Scan(&token.ID, &token.AccessToken, &token.RefreshToken, &token.ExpiresAt, &token.UserID)

	if err != nil {
		return nil, err
	}

	return token, nil
}
