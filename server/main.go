package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/otiai10/openaigo"
	"github.com/rs/cors"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"os"
)

var apiKey = os.Getenv("PLAYLISTIFY_CHATGPT_KEY")
var port = os.Getenv("PLAYLISTIFY_PORT")

type RequestBody struct {
	Msg string `json:"msg"`
}

type ChatServer struct {
	client *openaigo.Client
	db     *sql.DB
}

func (srv *ChatServer) Chat(w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*User)
	if user == nil {
		w.Write([]byte("no user"))
		w.WriteHeader(400)
		return
	}
	token, err := getTokenByUserID(srv.db, user.ID)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		token, err = RefreshSpotifyAccessToken(token)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(400)
			return
		}
		upsertToken(srv.db, token)
	}

	rb := &RequestBody{}
	err = json.NewDecoder(req.Body).Decode(rb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if rb == nil || rb.Msg == "" {
		w.WriteHeader(503)
		w.Write([]byte("no message"))
		return
	}

	chreq := openaigo.ChatCompletionRequestBody{
		Model:       "gpt-3.5-turbo",
		Temperature: 0.5,
		Messages: []openaigo.ChatMessage{
			{Role: "user", Content: "generate a playlist of 12 to 16 real and existing songs based on the following: `" + rb.Msg + "`\n layout your response in the following json template: `{ songs: {artist: string, title: string}[], description: string, title: string }`. in the description field please explain your logic for the playlist. give the playlist a really fun title based on your logic. do not give me any data outside the json. Make sure to include only tracks that have been officially released on Spotify and avoid any made-up or reimagined songs."},
		},
	}

	response, err := srv.client.Chat(context.Background(), chreq)
	if err != nil {
		w.WriteHeader(503)
		w.Write([]byte(err.Error()))
		return
	}
	b, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

func main() {
	// Initiate AWS Lambda handler
	// Finally, we call the `ListenAndServe` with the port
	// and a handler. `nil` tells it to use the default
	// router we've just set up.
	if port == "" {
		log.Fatalf("no port number provided: %s")
	}
	if apiKey == "" {
		log.Fatalf("no api key provided")
	}

	if clientSecret == "" || clientID == "" || redirectURI == "" {
		log.Fatalf("OS env keys not set (clientSecret: %s; clientID: %s; redirectURI: %s;)", clientSecret, clientID, redirectURI)
	}

	if JWTSecret == "" {
		log.Fatalf("OS env keys not set (JWTSecret: %s)", JWTSecret)
	}

	mux := http.NewServeMux()
	db, err := getConnection()

	if err != nil {
		log.Fatal("db error: " + err.Error())
	}

	log.Println("database connected")

	srv := ChatServer{
		client: openaigo.NewClient(apiKey),
		db:     db,
	}
	mux.Handle("/chat", srv.withUser(srv.Chat))
	mux.Handle("/me", srv.withUser(srv.Me))
	mux.Handle("/auth", srv.withUser(srv.Auth))
	mux.Handle("/spotify/api/", srv.withUser(srv.ForwardToSpotifyAPI))
	mux.HandleFunc("/spotify/callback", srv.AuthHandler)
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Println(request.URL.Path)
		http.Error(writer, "nothing here", 404)
	})

	handler := cors.Default().Handler(mux)
	fmt.Println("Serving Playlistify API on port " + port)
	http.ListenAndServe(":"+port, handler)

}
