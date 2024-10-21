package model

import "github.com/bwmarrin/discordgo"

// CommandHandler is a function type for handling commands
type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate) error

// Command represents a command with its handler
type Command struct {
	Name        string
	Description string
	Handler     CommandHandler
}
