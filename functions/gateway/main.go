package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"nasubaby-werewolf-hoster/cmd/command_init"
	"nasubaby-werewolf-hoster/cmd/commands"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	db_key := os.Getenv("DB_KEY")
	public_key := os.Getenv("PUBLIC_KEY")
	bot_token := os.Getenv("BOT_TOKEN")
	app_id := os.Getenv("APP_ID")
	if public_key == "" || bot_token == "" || app_id == "" {
		return Response{}, errors.New("Couldn't get enviroment variables")
	}

	pubkey_b, err := hex.DecodeString(public_key)

	if err != nil {
		return Response{}, errors.New("Couldn't decode the public key")
	}
	if request.Body == "" {
		if request.QueryStringParameters["command"] == "update" {
			command_init.Run(app_id, bot_token)
			return Response{
				StatusCode: 200,
				Body:       `{"log":"delete and update done"}`,
			}, nil
		}
		log.Print("400 No data")
		return Response{
			StatusCode: 400,
			Body:       `{"error":"No body data"}`,
		}, nil
	}

	var body []byte

	if request.IsBase64Encoded {
		body_b, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return Response{}, errors.New(fmt.Sprintf("Couldn't decode request body [%s]: %s", body, err))
		}
		body = body_b
	} else {
		body = []byte(request.Body)
	}
	pubkey := ed25519.PublicKey(pubkey_b)

	XSig, ok := request.Headers["x-signature-ed25519"]

	if !ok {
		log.Print("400 No Signature header")
		return Response{
			StatusCode: 400,
			Body:       `{"error": "Missing 'X-Signature-Ed25519' header"}`,
		}, nil
	}

	XSigTime, ok := request.Headers["x-signature-timestamp"]

	if !ok {
		log.Print("400 Missing Timestamp header")
		return Response{
			StatusCode: 400,
			Body:       `{"error": "Missing 'X-Signature-Timestamp' header"}`,
		}, nil
	}

	XSigB, err := hex.DecodeString(XSig)

	if err != nil {
		return Response{}, errors.New("Couldn't decode signature")
	}

	SignedData := []byte(XSigTime + string(body))

	if !ed25519.Verify(pubkey, SignedData, XSigB) {
		log.Print("401 Unauthorized")
		return Response{
			StatusCode: 401,
		}, nil
	} else {
		//authorized
		var inter discordgo.Interaction
		err := json.Unmarshal(body, &inter)

		if err != nil {
			log.Printf("Error decoding interaction: %s", err)
			return Response{
				StatusCode: 400,
			}, nil
		}

		switch {
		case inter.Type == 1:
			{
				log.Print("200 Type 1 Ping")
				return Response{
					StatusCode: 200,
					Body:       `{"type":1}`,
				}, nil
			}
		case inter.Type == 2:
			{
				var s *discordgo.Session
				s, err = discordgo.New("Bot " + bot_token)
				log.Printf("Application command [%s]", inter.ApplicationCommandData().Name)
				log.Printf(string(body))
				handler, ok := commands.Commands[inter.ApplicationCommandData().Name]
				if !ok {
					return Response{
						StatusCode: 404,
						Body:       `{"error": "Command not found"}`,
					}, nil
				}

				err := handler.Handler(s, &inter, db_key)

				if err != nil {
					log.Printf("Error in handler: %s", err)
					return Response{}, err
				} else {
					log.Printf("Successfully use command")
					if err != nil {
						return Response{}, err
					}

					return Response{
						StatusCode: 200,
						Body:       "Successfully use command",
					}, nil
				}
			}
		default:
			{
				return Response{
					StatusCode: 501,
				}, nil
			}
		}
	}

}

func main() {
	lambda.Start(Handler)
}
