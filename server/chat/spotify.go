package main

import (
	"encoding/base64"
	"encoding/json"
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

var authedRedir = os.Getenv("PLAYLISTIFY_AUTH_REDIR")

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

func init() {
	if clientSecret == "" || clientID == "" {
		log.Fatal("spotify tokens not set")
	}
}

func (srv *ChatServer) RefreshSpotifyToken(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := getUser(srv.db, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	token, err := getTokenByUserID(srv.db, user.ID)
	if err != nil {
		http.Error(w, "Token not found", http.StatusNotFound)
		return
	}

	formData := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {token.RefreshToken},
	}

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(formData.Encode()))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))))))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch new token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	token.AccessToken = tokenResponse.AccessToken
	token.ExpiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	token, err = upsertToken(srv.db, token)
	if err != nil {
		http.Error(w, "Failed to update token in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
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

func generateJWT(user *User, token *TokenResponse) (string, error) {
	claims := jwt.MapClaims{
		"user_id":      user.ID,
		"access_token": token.AccessToken,
		"exp":          token.ExpiresAt.Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (srv *ChatServer) AuthHandler(w http.ResponseWriter, r *http.Request) {
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

	user, err := upsertUser(srv.db, &User{
		SpotifyID:   profile.ID,
		DisplayName: profile.DisplayName,
		Email:       profile.Email,
	})

	if err != nil {
		log.Fatal(err)
	}

	// Save tokens in the database
	// ...

	jwtToken, err := generateJWT(user, tokens)
	if err != nil {
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		log.Println("unable to generate JWT", err)
		return
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(tokens.ExpiresIn) * time.Second)

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		Path:     "/",
		Domain:   "", // Set your domain here for all subdomains
		Expires:  expiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
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

	fmt.Println("STATUS?", http.StatusOK)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("ERROR OK", redirectURI)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ERROR 3:", err.Error())
			log.Println("Error reading response body: %v", err)
		}

		fmt.Println(string(body))
		return nil, fmt.Errorf("failed to fetch tokens, status code: %d", resp.StatusCode)
	}

	var tokenResponse TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		fmt.Println("ERROR 4:", err.Error())
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
	defer resp.Body.Close()

	var profile SpotifyUserProfile
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

type SpotifyUserProfile struct {
	DisplayName  string `json:"display_name"`
	Email        string `json:"email"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		Height int    `json:"height"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
	} `json:"images"`
	Product string `json:"product"`
	Type    string `json:"type"`
	URI     string `json:"uri"`
}
