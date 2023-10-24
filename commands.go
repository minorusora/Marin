package main

import (
	"github.com/bwmarrin/discordgo"
)

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commandData := i.ApplicationCommandData()

	// Check the name of the command
	switch commandData.Name {
	case "avatar":
		authorName := i.Member.User.Username
		authorImageURL := i.Member.User.AvatarURL("256") // You can specify the size you want

		// Create a new embed
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name: authorName,
			},
			Image: &discordgo.MessageEmbedImage{
				URL: authorImageURL,
			},
		}

		// Create an interaction response with the embed
		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		}

		// Send the response back to Discord
		s.InteractionRespond(i.Interaction, response)
	}
}
