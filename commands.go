package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func interactionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type == discordgo.InteractionApplicationCommand {
		switch interaction.ApplicationCommandData().Name {
		case "avatar":
			if len(interaction.ApplicationCommandData().Options) > 0 {
				option := interaction.ApplicationCommandData().Options[0]
				if option.Type == discordgo.ApplicationCommandOptionUser {
					userName := option.UserValue(session).Username
					userAvatar := option.UserValue(session).AvatarURL("256")
					currentTime := time.Now()

					embed := &discordgo.MessageEmbed{
						Title: userName,
						Footer: &discordgo.MessageEmbedFooter{
							IconURL: interaction.Member.AvatarURL("64"),
							Text:    interaction.Member.User.Username + " istedi.",
						},
						Image: &discordgo.MessageEmbedImage{
							URL: userAvatar,
						},
						Timestamp: currentTime.Format(time.RFC3339),
					}

					response := discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{embed},
						},
					}

					session.InteractionRespond(interaction.Interaction, &response)
					return
				}
			} else {
				currentTime := time.Now()
				embed := &discordgo.MessageEmbed{
					Title: interaction.Member.User.Username,
					Image: &discordgo.MessageEmbedImage{
						URL: interaction.Member.User.AvatarURL("256"),
					},
					Timestamp: currentTime.Format(time.RFC3339),
				}
				response := &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
					},
				}
				session.InteractionRespond(interaction.Interaction, response)
			}
		case "rolsec":
			roles, err := session.GuildRoles(interaction.GuildID)
			if err != nil {
				log.Println("Roller çekilemedi:", err)
				return
			}

			var options []discordgo.SelectMenuOption
			for _, role := range roles {
				options = append(options, discordgo.SelectMenuOption{
					Label: role.Name,
					Value: role.ID,
				})
			}

			selectMenu := discordgo.SelectMenu{
				CustomID:    "roleSelect",
				Placeholder: "Seç",
				Options:     options,
			}

			actionRow := discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{selectMenu},
			}

			response := discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content:    "Girişte üyelere verilecek rolü seçiniz:",
					Components: []discordgo.MessageComponent{actionRow},
				},
			}

			session.InteractionRespond(interaction.Interaction, &response)
		}
	} else if interaction.Type == discordgo.InteractionMessageComponent {
		customID := interaction.MessageComponentData().CustomID
		if customID == "roleSelect" {
			selectedRoleID := interaction.MessageComponentData().Values[0]

			selectedRole, err := session.State.Role(interaction.GuildID, selectedRoleID)
			if err != nil {
				log.Println("Error getting role:", err)
				return
			}

			responseText := selectedRole.Name + " rolü artık yeni bir üye katıldığında verilecek."

			response := discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: responseText,
				},
			}

			session.InteractionRespond(interaction.Interaction, &response)

			err = session.ChannelMessageDelete(interaction.ChannelID, interaction.Message.ID)
			if err != nil {
				return
			}

			db, err := sql.Open("mysql", dsn(dbname))
			if err != nil {
				log.Printf("%s, veritaban bağlantı hatası.\n", err)
				return
			}

			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM sunucuveri WHERE sunucu_id = ?", interaction.GuildID).Scan(&count)
			if err != nil {
				panic(err.Error())
			}
			if count > 0 {
				_, err := db.Exec("UPDATE sunucuveri SET rol = ? WHERE sunucu_id = ?", selectedRole.ID, interaction.GuildID)
				if err != nil {
					panic(err.Error())
				}
			} else {
				_, err := db.Exec("insert into sunucuveri(sunucu_id, rol) values (?, ?)", interaction.GuildID, selectedRole.ID)
				if err != nil {
					panic(err.Error())
				}
			}

			defer db.Close()
		}
	}
}

func ErrorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}
