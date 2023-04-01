package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
)

var JWTSecret string = os.Getenv("PLAYLISTIFY_JWT_SECRET")

func (srv *ChatServer) withUser(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			w.Write([]byte("no cookie found"))
			w.WriteHeader(400)
			return
		}
		if err == nil && cookie != nil {
			tokenString := cookie.Value

			jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(JWTSecret), nil
			})

			if err != nil {
				w.Write([]byte("no jwt: " + err.Error()))
				w.WriteHeader(400)
				return
			}

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

func (srv *ChatServer) Auth(w http.ResponseWriter, r *http.Request) {
	authEndpoint := "https://accounts.spotify.com/authorize"

	user := r.Context().Value("user").(*User)
	if user != nil {
		token, err := getTokenByUserID(srv.db, user.ID)
		if token == nil || err != nil {
			http.Redirect(w, r, authEndpoint, 302)
			return
		}

		if token.ExpiresAt.Before(time.Now()) {
			token, err = RefreshSpotifyAccessToken(token)
			if err != nil {
				http.Redirect(w, r, authEndpoint, 302)
				return
			}
			upsertToken(srv.db, token)
		}

		http.Redirect(w, r, "/", 302)
		return
	}

	http.Redirect(w, r, authEndpoint, 302)
	return
}

func (srv *ChatServer) Me(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	if user == nil {
		w.Write([]byte("no user"))
		w.WriteHeader(400)
	}

	token, err := getTokenByUserID(srv.db, user.ID)
	if err != nil || token == nil {
		w.Write([]byte("no access token"))
		w.WriteHeader(400)
		return
	}

	response := struct {
		AccessToken string    `json:"access_token"`
		ExpiresAt   time.Time `json:"expires_at"`
		*User
	}{
		AccessToken: token.AccessToken,
		ExpiresAt:   token.ExpiresAt,
		User:        user,
	}

	// Marshal the user model as JSON
	userJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Set the response content type and write the user JSON to the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}

func getUser(db *sql.DB, userID string) (*User, error) {
	query := `SELECT id, spotify_id, display_name, created_at FROM users WHERE id = $1`
	row := db.QueryRow(query, userID)

	var user User
	err := row.Scan(&user.ID, &user.SpotifyID, &user.DisplayName, &user.CreatedAt)

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
		SELECT access_token, refresh_token, expires_at, user_id
		FROM tokens
		WHERE user_id = $1;
	`

	row := db.QueryRow(query, userID)

	var token Token
	err := row.Scan(&token.AccessToken, &token.RefreshToken, &token.ExpiresAt, &token.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no token found for user ID: %s", userID)
		}
		return nil, err
	}

	return &token, nil
}
