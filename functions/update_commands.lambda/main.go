package main

import (
	"context"
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

func DeleteCommands(app_id string, disc *discordgo.Session) error {
	cmds, err := disc.ApplicationCommands(app_id, "")

	if err != nil {
		log.Printf("Error fetching commands")
		return err
	}

	for j := range cmds {
		if err := disc.ApplicationCommandDelete(app_id, "", cmds[j].ID); err != nil {
			log.Printf("Error deleting command [%s]: %s", cmds[j].Name, err)
			return err
		} else {
			log.Printf("Deleted %s\n", cmds[j].Name)
			return nil
		}
	}
	return nil

}

// Handler function Using AWS Lambda Proxy Request
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	app_id := os.Getenv("APP_ID")
	bot_token := os.Getenv("BOT_TOKEN")

	if bot_token == "" || app_id == "" {
		return Response{StatusCode: 500, Body: "couldn't get environment variables"}, errors.New("couldn't get environment variables")
	}

	disc, err := discordgo.New("Bot " + bot_token)

	if err != nil {
		log.Printf("Error creating discord client: %s\n", err.Error())
		return Response{StatusCode: 500, Body: "Error creating discord client"}, err
	}

	err = disc.Open()

	if err != nil {
		log.Printf("Error opening a connection to discord: %s\n", err.Error())
		return Response{StatusCode: 500, Body: "Error opening a connection to discord"}, err
	}
	//Delete all commands
	err = DeleteCommands(app_id, disc)

	if err != nil {
		return Response{StatusCode: 500, Body: "Error deleting commands"}, err
	}

	//Create new commands
	for j := range commands.Commands {
		cmd := discordgo.ApplicationCommand{
			Name:        commands.Commands[j].Name,
			Description: commands.Commands[j].Description,
		}
		_, err := disc.ApplicationCommandCreate(app_id, "", &cmd)
		if err != nil {
			log.Printf("Error creating command [%s]: %s", commands.Commands[j].Name, err)
			return Response{StatusCode: 500, Body: "Error creating command"}, err
		}
	}

	return Response{
		StatusCode: 200,
		Body:       `{"message": "Commands updated successfully"}`,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
