package command_init

import (
	"log"
	"nasubaby-werewolf-hoster/cmd/commands"

	"github.com/bwmarrin/discordgo"
)

func Delete(app_id string, disc *discordgo.Session) {
	cmds, err := disc.ApplicationCommands(app_id, "")

	if err != nil {
		log.Fatal("Error fetching commands")
	}

	for j := range cmds {
		if err := disc.ApplicationCommandDelete(app_id, "", cmds[j].ID); err != nil {
			log.Fatalf("Error deleting command [%s]: %s", cmds[j].Name, err)
		} else {
			log.Printf("Deleted %s\n", cmds[j].Name)
		}
	}

}

func Run(app_id string, bot_token string) {
	disc, err := discordgo.New("Bot " + bot_token)

	if err != nil {
		log.Fatalf("Error creating discord client: %s\n", err.Error())
	}

	err = disc.Open()

	if err != nil {
		log.Fatalf("Error opening a connection to discord: %s\n", err.Error())
	}

	Delete(app_id, disc)

	for j := range commands.Commands {
		cmd := commands.Commands[j].Command
		_, err := disc.ApplicationCommandCreate(app_id, "", &cmd)
		if err != nil {
			log.Fatalf("Error uploading command [%s]: %s", cmd.Name, err)
		} else {
			log.Printf("Uploaded command [%s]\n", cmd.Name)
		}
	}

}
