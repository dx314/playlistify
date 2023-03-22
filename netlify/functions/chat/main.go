package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/otiai10/openaigo"
	"os"
)

var client = openaigo.NewClient(os.Getenv("CHATGPT_KEY"))

func Handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	msg := request.QueryStringParameters["msg"]

	chreq := openaigo.ChatCompletionRequestBody{
		Model: "gpt-3.5-turbo",
		Messages: []openaigo.ChatMessage{
			{Role: "user", Content: msg},
		},
	}

	response, err := client.Chat(context.Background(), chreq)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 503,
			Body:       err.Error(),
		}, nil
	}
	b, _ := json.Marshal(response)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(b),
	}, nil
}

func main() {
	// Initiate AWS Lambda handler
	lambda.Start(Handler)
}
