package handler

import "github.com/bwmarrin/discordgo"

// HelloHandler handles the /hello command
func HelloHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hello, world!",
		},
	}
	return s.InteractionRespond(i.Interaction, response)
}
