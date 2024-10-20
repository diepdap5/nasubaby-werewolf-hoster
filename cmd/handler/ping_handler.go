package handler

import "github.com/bwmarrin/discordgo"

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
