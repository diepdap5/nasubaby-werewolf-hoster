package handler

import (
	"github.com/bwmarrin/discordgo"
)

func PingHandler(s *discordgo.Session, interaction *discordgo.Interaction, db_key string) error {
	s.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
	return nil
}
