package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var clientSecret = os.Getenv("PLAYLISTIFY_SPOTIFY_SECRET")
var clientID = os.Getenv("PLAYLISTIFY_SPOTIFY_ID")
var redirectURI = os.Getenv("PLAYLISTIFY_SPOTIFY_REDIR")

var scopes = []string{
	"user-read-private",
	"playlist-modify-private",
	"playlist-modify-public",
	"playlist-read-private",
	"playlist-read-collaborative",
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func RefreshSpotifyAccessToken(token *Token) (*Token, error) {
	// Create a new HTTP client
	if token == nil || token.RefreshToken == "" {
		return nil, errors.New("no token or no refresh token")
	}

	client := &http.Client{}

	// Create a new URL for the Spotify token endpoint
	tokenUrl := "https://accounts.spotify.com/api/token"

	// Create a new request to the Spotify token endpoint
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", token.RefreshToken)
	req, err := http.NewRequest("POST", tokenUrl, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}

	// Add the Authorization header with the client credentials
	clientCreds := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", clientCreds))

	// Set the Content-Type header to "application/x-www-form-urlencoded"
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request to the Spotify token endpoint
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Decode the response body into a map
	var body map[string]interface{}

	bodystr, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(bodystr, &body)
	if err != nil {
		return nil, err
	}

	// Get the access token from the response body
	accessToken, ok := body["access_token"].(string)
	if !ok {
		return nil, fmt.Errorf("could not get access token from response")
	}

	refreshToken, ok := body["refresh_token"].(string)

	token.AccessToken = accessToken
	if ok && refreshToken != "" {
		token.RefreshToken = refreshToken
	}

	// Get the expiration time from the response body
	expirySeconds, ok := body["expires_in"].(float64)
	if !ok {
		return nil, fmt.Errorf("could not get token expiry from response")
	}
	token.ExpiresAt = time.Now().Add(time.Duration(expirySeconds) * time.Second)

	// Return the access token and the expiry time
	return token, nil
}

// TokenResponse is the structure for the token response from the Spotify API
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"-"`
	RefreshToken string    `json:"refresh_token"`
	Scope        string    `json:"scope"`
	TokenType    string    `json:"token_type"`
}

func generateJWT(user *User, token *Token) (string, error) {
	claims := jwt.MapClaims{
		"user_id":      user.ID,
		"access_token": token.AccessToken,
		"exp":          token.ExpiresAt.Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (srv *ChatServer) AuthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting auth token")

	authCode := r.URL.Query().Get("code")
	if authCode == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	tokens, err := fetchTokens(authCode)
	if err != nil {
		http.Error(w, "Failed to fetch tokens", http.StatusInternalServerError)
		log.Println("unable to fetch tokens", err)
		return
	}

	profile, err := fetchUserProfile(tokens.AccessToken)
	if err != nil {
		http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
		log.Println("unable to fetch user profile", err)
		return
	}

	// Perform logic to find or create the user in the database using profile.ID
	// ...
	user := &User{
		SpotifyID:   profile.ID,
		DisplayName: profile.DisplayName,
	}

	err = upsertUser(srv.db, user)

	if err != nil {
		log.Fatal(err)
	}

	token := &Token{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		UserID:       user.ID,
	}

	token.ExpiresAt = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
	err = upsertToken(srv.db, token)
	if err != nil {
		http.Error(w, "Failed to write tokens to database", http.StatusInternalServerError)
		log.Println("unable to write tokens to database", err)
		return
	}
	// Save tokens in the database
	// ...

	jwtToken, err := generateJWT(user, token)
	if err != nil {
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		log.Println("unable to generate JWT", err)
		return
	}

	expiration := time.Now().Add(30 * 24 * time.Hour)
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		Path:     "/",
		Domain:   "plailist.app", // Set your domain here for all subdomains
		Expires:  expiration,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}

func fetchTokens(authCode string) (*TokenResponse, error) {
	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", authCode)
	data.Add("redirect_uri", redirectURI)
	data.Add("scope", url.QueryEscape(strings.Join(scopes, " ")))

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", ioutil.NopCloser(strings.NewReader(data.Encode())))
	if err != nil {

		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failure to fetch token: %s", string(body))
			log.Println("Error reading response body: %v", err)
		}

		return nil, fmt.Errorf("failed to fetch tokens, status code: %d", resp.StatusCode)
	}

	var tokenResponse TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func fetchUserProfile(accessToken string) (*SpotifyUserProfile, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error:", err)
	}

	resp.Body.Close()

	var profile SpotifyUserProfile
	err = json.Unmarshal(body, &profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

type SpotifyUserProfile struct {
	Country         string `json:"country,omitempty"`
	DisplayName     string `json:"display_name,omitempty"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled,omitempty"`
		FilterLocked  bool `json:"filter_locked,omitempty"`
	} `json:"explicit_content,omitempty"`
	ExternalUrls struct {
		Spotify string `json:"spotify,omitempty"`
	} `json:"external_urls,omitempty"`
	Followers struct {
		Href  any `json:"href,omitempty"`
		Total int `json:"total,omitempty"`
	} `json:"followers,omitempty"`
	Href    string `json:"href,omitempty"`
	ID      string `json:"id,omitempty"`
	Images  []any  `json:"images,omitempty"`
	Product string `json:"product,omitempty"`
	Type    string `json:"type,omitempty"`
	URI     string `json:"uri,omitempty"`
}
