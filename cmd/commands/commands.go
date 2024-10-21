package commands

import (
	"nasubaby-werewolf-hoster/cmd/handler"
	"nasubaby-werewolf-hoster/cmd/model"
)

// Commands is a map of command names to their handlers
var Commands = map[string]model.Command{
	"ping": {
		Name:        "ping",
		Description: "A simple ping command",
		Handler:     handler.PingHandler,
	},
	"hello": {
		Name:        "hello",
		Description: "Responds with 'Hello, world!'",
		Handler:     handler.HelloHandler,
	},
	"listgamebase": {
		Name:        "listgamebase",
		Description: "List all game bases",
		Handler:     handler.ListGameBaseHandler,
	},
}
