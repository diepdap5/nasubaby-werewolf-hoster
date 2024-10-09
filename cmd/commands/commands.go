package commands

import (
	"github.com/bwmarrin/discordgo"
)

// CommandHandler is a function type for handling commands
type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate, dbKey string) error

// Command represents a command with its handler
type Command struct {
	Name    string
	Handler CommandHandler
}

// Commands is a map of command names to their handlers
var Commands = map[string]Command{
	"ping": {
		Name:    "ping",
		Handler: PingHandler,
	},
}

// PingHandler handles the /ping command
func PingHandler(s *discordgo.Session, i *discordgo.InteractionCreate, dbKey string) error {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	}
	return s.InteractionRespond(i.Interaction, response)
}
