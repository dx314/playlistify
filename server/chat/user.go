package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
)

func (srv *ChatServer) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")
		if err == nil && cookie != nil {
			tokenString := cookie.Value

			jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err == nil && jwtToken.Valid {
				claims := jwtToken.Claims.(jwt.MapClaims)
				userID := claims["user_id"].(string)

				user, err := getUser(srv.db, userID)
				if err == nil {
					ctx := context.WithValue(r.Context(), "user", user)
					r = r.WithContext(ctx)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func getUser(db *sql.DB, userID string) (*User, error) {
	query := `SELECT id, spotify_id, display_name, email, created_at FROM users WHERE id = $1`
	row := db.QueryRow(query, userID)

	var user User
	err := row.Scan(&user.ID, &user.SpotifyID, &user.DisplayName, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func getTokenByUserID(db *sql.DB, userID string) (*Token, error) {
	query := `
		SELECT id, access_token, refresh_token, expires_in, user_id
		FROM tokens
		WHERE user_id = $1;
	`

	row := db.QueryRow(query, userID)

	var token Token
	err := row.Scan(&token.ID, &token.AccessToken, &token.RefreshToken, &token.ExpiresAt, &token.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no token found for user ID: %s", userID)
		}
		return nil, err
	}

	return &token, nil
}
