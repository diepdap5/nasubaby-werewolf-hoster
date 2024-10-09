package commands

import (
	"babynasu-queen/cmd/handler"
	"log"

	"github.com/bwmarrin/discordgo"
)

type (
	HandlerSig      = func(*discordgo.Session, *discordgo.Interaction, string) error
	ContinuationSig = func(*discordgo.Interaction) error
)

// checks for inconsistencies in Commands
func init() {
	for key := range Commands {
		if Commands[key].Command.Name != key {
			log.Fatalf("Error: Key [%s] in Commands doesn't equal the Name property [%s]\n", key, Commands[key].Command.Name)
		}
		if Commands[key].Handler == nil {
			log.Fatalf("Error: Handler for command [%s] is nil \n", key)
		}
	}
}

type Command struct {
	Command      discordgo.ApplicationCommand
	Handler      HandlerSig
	Continuation ContinuationSig
}

var Commands = map[string]Command{
	"ping": {
		Command: discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Send ping to get pong",
		},
		Handler: handler.PingHandler,
	},
	"gamestart": {
		Command: discordgo.ApplicationCommand{
			Name:        "gamestart",
			Description: "Start game werewolf!",
		},
		Handler: handler.StartGameHandler,
	},
}
