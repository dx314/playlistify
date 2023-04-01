package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

import (
	"fmt"
)

func requestSpotifyAPI(method string, path string, accessToken string, queryParams url.Values, body io.Reader) (*http.Response, error) {
	// Define the base URL of the Spotify API
	baseUrl := "https://api.spotify.com/v1"

	// Define the target URL by appending the request path to the base URL
	targetUrl := fmt.Sprintf("%s%s", baseUrl, path)

	// Create a new HTTP request to the target URL
	req, err := http.NewRequest(method, targetUrl, body)
	if err != nil {
		return nil, err
	}

	bearer := fmt.Sprintf("Bearer %s", accessToken)
	req.URL.RawQuery = queryParams.Encode()

	// Add the access token to the request headers
	req.Header.Set("Authorization", bearer)

	// Make the request to the Spotify API
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (srv *ChatServer) ForwardToSpotifyAPI(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request context
	user, ok := r.Context().Value("user").(*User)
	if !ok || user == nil {
		http.Error(w, "user ID not found in request context", http.StatusInternalServerError)
		return
	}

	// Get the access token for the user ID
	token, err := getTokenByUserID(srv.db, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the access token has expired

	if time.Now().Before(token.ExpiresAt.Add(-5 * time.Minute)) {
		// Access token has expired, refresh it
		token, err = RefreshSpotifyAccessToken(token)
		if err != nil {
			http.Error(w, "Reset token failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		upsertToken(srv.db, token)
	}

	// Get the path of the request
	path := strings.TrimPrefix(r.URL.Path, "/spotify/api")

	queryParams := r.URL.Query()
	// Make the request to the Spotify API using the requestSpotifyAPI function
	resp, err := requestSpotifyAPI(r.Method, path, token.AccessToken, queryParams, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers to the response writer
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set the status code of the response writer to match the status code of the response
	w.WriteHeader(resp.StatusCode)

	// Copy the response body to the response writer
	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
