package commands

import (
	"github.com/bwmarrin/discordgo"
)

// CommandHandler is a function type for handling commands
type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate) error

// Command represents a command with its handler
type Command struct {
	Name        string
	Description string
	Handler     CommandHandler
}

// Commands is a map of command names to their handlers
var Commands = map[string]Command{
	"ping": {
		Name:        "ping",
		Description: "A simple ping command",
		Handler:     PingHandler,
	},
	"hello": {
		Name:        "hello",
		Description: "Responds with 'Hello, world!'",
		Handler:     HelloHandler,
	},
}

// PingHandler handles the /ping command
func PingHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	}
	return s.InteractionRespond(i.Interaction, response)
}

// HelloHandler handles the /hello command
func HelloHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hello, world!",
		},
	}
	return s.InteractionRespond(i.Interaction, response)
}
