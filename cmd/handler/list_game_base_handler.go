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
	// Get database connection
	db, err := db.GetDBConnection()
	if err != nil {
		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to get game bases due to database connection error",
			},
		}
		return s.InteractionRespond(i.Interaction, response)
	}
	// Get all game bases
	gamebase, err := repository.GetAll(db)
	if err != nil {
		log.Println(err)
		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to get game bases",
			},
		}
		return s.InteractionRespond(i.Interaction, response)
	}
	content := "Game Bases:\n"
	for _, gb := range gamebase {
		rolesList := make([]string, len(gb.RolesList))
		for i, role := range gb.RolesList {
			rolesList[i] = role.Name // Assuming model.Role has a Name field
		}
		content += strconv.Itoa(gb.ID+1) + ": " + strconv.Itoa(gb.RoleCount) + ": " + "\n" + "- " + strings.Join(rolesList, "\n- ") + "\n"
	}
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}
	return s.InteractionRespond(i.Interaction, response)
}
