package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/otiai10/openaigo"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

var apiKey = os.Getenv("PLAYLISTIFY_CHATGPT_KEY")
var port = os.Getenv("PLAYLISTIFY_PORT")

type RequestBody struct {
	Msg string `json:"msg"`
}

type ChatServer struct {
	client *openaigo.Client
}

func (cs *ChatServer) Chat(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	rb := &RequestBody{}
	err := json.NewDecoder(req.Body).Decode(rb)
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
			{Role: "user", Content: "generate a playlist of 12 songs based on the following: `" + rb.Msg + "`\n layout your response in the following json template: `{ songs: {artist: string, title: string}[], description: string, title: string }`. in the description field please explain your logic for the playlist. give the playlist a really fun title based on your logic. do not give me any data outside the json. make sure all the songs actually exist."},
		},
	}

	response, err := cs.client.Chat(context.Background(), chreq)
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
		log.Fatal("no port number provided")
	}
	if apiKey == "" {
		log.Fatal("no api key provided")
	}
	mux := http.NewServeMux()

	srv := ChatServer{
		client: openaigo.NewClient(apiKey),
	}
	mux.HandleFunc("/chat", srv.Chat)
	handler := cors.Default().Handler(mux)
	fmt.Println("Serving Playlistify API on port " + port)
	http.ListenAndServe(":"+port, handler)

}
