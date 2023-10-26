package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func interactionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type == discordgo.InteractionApplicationCommand {
		currentTime := time.Now()
		switch interaction.ApplicationCommandData().Name {
		case "yardım":
			embed := &discordgo.MessageEmbed{
				Title:  "YARDIM KOMUTLARI",
				Fields: yardimEmbed,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Marin Geliştirici Ekibi",
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: session.State.User.AvatarURL("128"),
				},
				Timestamp: currentTime.Format(time.RFC3339),
			}

			embedGonder(session, interaction, embed)
		case "avatar":
			if len(interaction.ApplicationCommandData().Options) > 0 {
				option := interaction.ApplicationCommandData().Options[0]
				if option.Type == discordgo.ApplicationCommandOptionUser {
					userName := option.UserValue(session).Username
					userAvatar := option.UserValue(session).AvatarURL("256")

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

					embedGonder(session, interaction, embed)
					return
				}
			} else {
				embed := &discordgo.MessageEmbed{
					Title: interaction.Member.User.Username,
					Image: &discordgo.MessageEmbedImage{
						URL: interaction.Member.User.AvatarURL("256"),
					},
					Timestamp: currentTime.Format(time.RFC3339),
				}
				embedGonder(session, interaction, embed)
			}
		case "girisayarla":
			if len(interaction.ApplicationCommandData().Options) > 0 {
				option := interaction.ApplicationCommandData().Options[0]
				option1 := interaction.ApplicationCommandData().Options[1]
				kanalID := ""
				str := ""

				if option.Type == discordgo.ApplicationCommandOptionChannel {
					kanalID = option.ChannelValue(session).ID
				} else if option.Type == discordgo.ApplicationCommandOptionString {
					kanalID = option.StringValue()
				}

				if option1.Type == discordgo.ApplicationCommandOptionString {
					str = option1.StringValue()
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
					_, err := db.Exec("UPDATE sunucuveri SET kanal_id = ?, giris_mesaj = ? WHERE sunucu_id = ?", kanalID, str, interaction.GuildID)
					if err != nil {
						panic(err.Error())
					}
				} else {
					_, err := db.Exec("insert into sunucuveri(sunucu_id, kanal_id, giris_mesaj) values (?, ?, ?)", interaction.GuildID, kanalID, str)
					if err != nil {
						panic(err.Error())
					}
				}

				defer db.Close()

				embed := &discordgo.MessageEmbed{
					Description: "<#" + kanalID + ">" + " kanalına gönderilecek mesaj: " + str,
				}
				embedGonder(session, interaction, embed)
			}
		case "kanalolustur":
			category := interaction.ApplicationCommandData().Options[0].ChannelValue(session).ID
			channelName := interaction.ApplicationCommandData().Options[1].StringValue()

			newChannel, err := createChannel(session, interaction.GuildID, channelName)
			if err != nil {
				return
			}

			err = moveChannelToCategory(session, newChannel.ID, category)
			if err != nil {
				return
			}
			embed := &discordgo.MessageEmbed{
				Description: "<#" + category + ">" + " kategorisine " + newChannel.Mention() + " kanalı oluşturuldu.",
				Timestamp:   currentTime.Format(time.RFC3339),
			}
			embedGonder(session, interaction, embed)
		case "kanalsil":
			kanal := interaction.ApplicationCommandData().Options[0].ChannelValue(session).ID
			_, err := session.ChannelDelete(kanal)
			if err != nil {
				return
			}
			embed := &discordgo.MessageEmbed{
				Description: kanal + " ID kanal silindi.",
				Timestamp:   currentTime.Format(time.RFC3339),
			}
			embedGonder(session, interaction, embed)
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

			embed := &discordgo.MessageEmbed{
				Description: selectedRole.Name + " rolü artık yeni bir üye katıldığında verilecek.",
			}

			embedGonder(session, interaction, embed)

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

func embedGonder(session *discordgo.Session, interaction *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}
	session.InteractionRespond(interaction.Interaction, response)
}

func createChannel(session *discordgo.Session, guildID, channelName string) (*discordgo.Channel, error) {
	channel, err := session.GuildChannelCreate(guildID, channelName, discordgo.ChannelTypeGuildText)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func moveChannelToCategory(session *discordgo.Session, channelID, categoryID string) error {
	editData := &discordgo.ChannelEdit{
		ParentID: categoryID,
	}
	_, err := session.ChannelEditComplex(channelID, editData)
	if err != nil {
		return err
	}
	return nil
}
