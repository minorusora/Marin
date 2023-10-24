package main

import (
	"github.com/bwmarrin/discordgo"
)

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commandData := i.ApplicationCommandData()

	switch commandData.Name {
	case "avatar":
		authorName := i.Member.User.Username
		authorImageURL := i.Member.User.AvatarURL("256")

		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name: authorName,
			},
			Image: &discordgo.MessageEmbedImage{
				URL: authorImageURL,
			},
		}

		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		}

		s.InteractionRespond(i.Interaction, response)
	}
}
