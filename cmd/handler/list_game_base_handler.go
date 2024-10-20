package handler

import (
	"log"
	"nasubaby-werewolf-hoster/cmd/db"
	"nasubaby-werewolf-hoster/cmd/repository"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// HelloHandler handles the /hello command
func ListGameBaseHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	db, err := db.GetDBConnection()
	if err != nil {
		return err
	}
	// Get all game bases
	gamebase, err := repository.GetGameBases(db)
	if err != nil {
		log.Println(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to get game bases",
			},
		})
		return err
	}
	content := "Game Bases:\n"
	for _, gb := range gamebase {
		content += strconv.Itoa(gb.ID+1) + ": " + strconv.Itoa(gb.RoleCount) + ": " + "\n" + "- " + strings.Join(gb.RolesList, "\n- ") + "\n"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	return nil
}
