package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"nasubaby-werewolf-hoster/cmd/commands"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
type Response events.APIGatewayProxyResponse

// Handler function Using AWS Lambda Proxy Request
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	public_key := os.Getenv("PUBLIC_KEY")
	bot_token := os.Getenv("BOT_TOKEN")
	app_id := os.Getenv("APP_ID")
	if public_key == "" || bot_token == "" || app_id == "" {
		return Response{StatusCode: 500, Body: "couldn't get environment variables"}, errors.New("couldn't get environment variables")
	}

	pubkey_b, err := hex.DecodeString(public_key)
	if err != nil {
		return Response{StatusCode: 500, Body: "couldn't decode the public key"}, errors.New("couldn't decode the public key")
	}
	pubkey := ed25519.PublicKey(pubkey_b)

	XSig, ok := request.Headers["x-signature-ed25519"]
	if !ok {
		log.Print("400 No Signature header")
		log.Print(request.Headers)
		return Response{
			StatusCode: 400,
			Body:       `{"error": "Missing 'x-signature-ed25519' header"}`,
		}, nil
	}

	XTime, ok := request.Headers["x-signature-timestamp"]
	if !ok {
		log.Print("400 No Timestamp header")
		return Response{
			StatusCode: 400,
			Body:       `{"error": "Missing 'x-signature-timestamp' header"}`,
		}, nil
	}

	body := request.Body
	SignedData := []byte(XTime + body)
	XSigB, err := hex.DecodeString(XSig)
	if err != nil {
		log.Print("401 Unauthorized - Invalid Signature")
		return Response{
			StatusCode: 401,
			Body:       `{"error": "Invalid 'x-signature-ed25519' header"}`,
		}, nil
	}

	if !ed25519.Verify(pubkey, SignedData, XSigB) {
		log.Print("401 Unauthorized")
		return Response{
			StatusCode: 401,
			Body:       `{"error": "Unauthorized"}`,
		}, nil
	}

	// Authorized
	var inter discordgo.InteractionCreate
	err = json.Unmarshal([]byte(body), &inter)
	if err != nil {
		log.Printf("Error decoding interaction: %s", err)
		return Response{
			StatusCode: 400,
			Body:       `{"error": "Invalid request body"}`,
		}, nil
	}

	switch inter.Type {
	case discordgo.InteractionPing:
		log.Print("200 Type 1 Ping")
		return Response{
			StatusCode: 200,
			Body:       `{"type":1}`,
		}, nil
	case discordgo.InteractionApplicationCommand:
		s, err := discordgo.New("Bot " + bot_token)
		if err != nil {
			log.Printf("Error creating Discord session: %s", err)
			return Response{
				StatusCode: 500,
				Body:       `{"error": "Internal server error"}`,
			}, nil
		}

		log.Printf("Application command [%s]", inter.ApplicationCommandData().Name)

		handler, ok := commands.Commands[inter.ApplicationCommandData().Name]
		if !ok {
			return Response{
				StatusCode: 404,
				Body:       `{"error": "Command not found"}`,
			}, nil
		}

		err = handler.Handler(s, &inter)
		if err != nil {
			log.Printf("Error in handler: %s", err)
			return Response{
				StatusCode: 500,
				Body:       `{"error": "Internal server error"}`,
			}, nil
		}

		// Add a successful response if needed
		return Response{
			StatusCode: 200,
			Body:       `{"message": "Command executed successfully"}`,
		}, nil
	default:
		return Response{
			StatusCode: 400,
			Body:       `{"error": "Unknown interaction type"}`,
		}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
