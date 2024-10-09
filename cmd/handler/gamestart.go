package handler

import (
	"babynasu-queen/cmd/repository"
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
)

func StartGameHandler(s *discordgo.Session, interaction *discordgo.Interaction, db_key string) error {
	gamebase, err := repository.GetByID(db_key, 2)
	if err != nil {
		log.Println(err)
		s.InteractionRespond(interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Start failed",
			},
		})
		return err
	}
	err = CreateRoleAndCreateChannel(gamebase, s, interaction)
	if err != nil {
		log.Println(err)
		s.InteractionRespond(interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Start failed",
			},
		})
		return err
	}
	log.Println("Start OK")
	s.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Start OK",
		},
	})
	return nil

}

func CreateRoleAndCreateChannel(g *repository.GameBase, s *discordgo.Session, interaction *discordgo.Interaction) (err error) {
	for _, role := range unique(g.RolesList) {
		created_role, err := s.GuildRoleCreate(interaction.GuildID, &discordgo.RoleParams{
			Name: role,
		})
		if err != nil {
			return errors.New("Error GuildRoleCreate")
		}
		if created_role.Name != "dân làng" {
			// Create Text Channel
			_, errText := s.GuildChannelCreateComplex(interaction.GuildID, discordgo.GuildChannelCreateData{
				Name: created_role.Name,
				Type: discordgo.ChannelTypeGuildText,
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:    created_role.ID,
						Allow: discordgo.PermissionAllText,
					},
				},
			})
			if errText != nil {
				return errors.New("Error create ChannelTypeGuildText")
			}
			// Create Voice Channel
			_, errVoice := s.GuildChannelCreateComplex(interaction.GuildID, discordgo.GuildChannelCreateData{
				Name: created_role.Name,
				Type: discordgo.ChannelTypeGuildVoice,
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:    created_role.ID,
						Allow: discordgo.PermissionAllVoice,
					},
				},
			})
			if errVoice != nil {
				return errors.New("Error create ChannelTypeGuildVoice")
			}
		}
	}
	return nil
}

func unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}
