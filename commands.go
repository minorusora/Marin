package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
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
		case "param":
			mesaj := fmt.Sprintf("%s, %s MC sahibisiniz! :money_with_wings:", interaction.Member.Mention(), formatNumber(paraCek(interaction.Member.User.ID)))
			embed := &discordgo.MessageEmbed{
				Description: mesaj,
				Timestamp:   currentTime.Format(time.RFC3339),
			}
			embedGonder(session, interaction, embed)
		case "hayvanfiyatları":
			embed := &discordgo.MessageEmbed{
				Title:       "Havan Fiyatları",
				Description: "İnek: `25.000 MC`\nKoyun: `12.500 MC`\nTavuk: `5.000 MC`",
				Timestamp:   currentTime.Format(time.RFC3339),
			}
			embedGonder(session, interaction, embed)
		case "seviyem":
			mesaj := fmt.Sprintf("%s, %s. seviyedesiniz. (%s/%s) :ringed_planet:", interaction.Member.Mention(), formatNumber(levelKontrol(interaction.Member.User.ID, interaction.GuildID)), formatNumber(xpCheck(interaction.Member.User.ID, interaction.GuildID)), formatNumber(levelKontrol(interaction.Member.User.ID, interaction.GuildID)*1024))
			embed := &discordgo.MessageEmbed{
				Description: mesaj,
				Timestamp:   currentTime.Format(time.RFC3339),
			}
			embedGonder(session, interaction, embed)
		case "çiftliğim":
			userID := interaction.Member.User.ID
			str := fmt.Sprintf("Çiftlik Seviyesi: `%s` | Hayvan Seviyesi: `%s`\n \nİnek Sayısı: `%s Adet` - Gelir: `%s MC`", formatNumber(ciftlikSeviye(userID)), formatNumber(hayvanSeviye(userID)), formatNumber(inekGet_Count(userID)), formatNumber(hayvanSeviye(userID)*10))
			str1 := fmt.Sprintf(str+"\nKoyun Sayısı: `%s Adet` - Gelir: `%s MC`\nTavuk Sayısı: `%s Adet` - Gelir: `%s MC`\n ", formatNumber(koyunGet_Count(userID)), formatNumber(hayvanSeviye(userID)*6), formatNumber(tavukGet_Count(userID)), formatNumber(hayvanSeviye(userID)*3))
			str2 := fmt.Sprintf(str1+"\nBuğday Sayısı: `%s/%s` - Buğday Tohumu Sayısı: `%s Adet`\nHavuç Sayısı: `%s/%s` - Havuç Tohumu Sayısı: `%s Adet`\n", formatNumber(mahsulSayi(1, userID)), formatNumber(ciftlikSeviye(userID)*25), formatNumber(tohumSayisi(1, userID)), formatNumber(mahsulSayi(2, userID)), formatNumber(ciftlikSeviye(userID)*25), formatNumber(tohumSayisi(2, userID)))
			mesaj := fmt.Sprintf(str2 + "\n **BİLGİLENDİRME**\n \n`Hayvan gelirleri çiftlik seviyesine bağlı olarak değişmez.`\n`Çiftlik seviyesi mahsül stoğunu arttırır.`\n \n`Gelirler 3 saatte bir verilir, toplam değil adet fiyatıdır.`\n \n**MAHSÜLLER**\n \n`Buğday tohumları 2 saatte bir yetişir.`\n`Havuç tohumları 1 saatte bir yetişir.`")
			embed := &discordgo.MessageEmbed{
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://cdn.discordapp.com/attachments/1169252458975989844/1173424309335769169/png-transparent-barn-farm-building-icon-warehouse-miscellaneous-building-cartoon-thumbnail-removebg-preview.png?ex=6563e78c&is=6551728c&hm=eb660c6cddaf42741a5ca6474c84ecb5d6f27a7a6d2b4ff0b4fcce4d82ef5fa9&",
				},
				Title:       interaction.Member.User.Username + " Çiftliği",
				Description: mesaj,
				Timestamp:   currentTime.Format(time.RFC3339),
			}
			embedGonder(session, interaction, embed)
		case "çiftlikseviye":
			kullaniciParasi := int64(paraCek(interaction.Member.User.ID))
			var seviye int = ciftlikSeviye(interaction.Member.User.ID)
			kalanPara := int(kullaniciParasi) - (seviye * 5000)
			if seviye*5000 <= int(kullaniciParasi) {
				paraKayit(session, interaction.Member.User.ID, int64(kalanPara))
				ciftlikSeviyeYukselt(interaction.Member.User.ID)
				mesaj := fmt.Sprintf("Çiftliğinizin, %s MC karşılığında seviyesi yükseltildi, şu an %d seviye çiftliğe sahipsiniz. :star2:", formatNumber(ciftlikSeviye(interaction.Member.User.ID)*5000), ciftlikSeviye(interaction.Member.User.ID))
				embed := &discordgo.MessageEmbed{
					Description: mesaj,
					Timestamp:   currentTime.Format(time.RFC3339),
				}
				embedGonder(session, interaction, embed)
			} else {
				embed := &discordgo.MessageEmbed{
					Description: fmt.Sprintf("Yeterli paranız yok. Toplam fiyat: %s MC :confused:", formatNumber(ciftlikSeviye(interaction.Member.User.ID)*5000)),
				}
				embedGonder(session, interaction, embed)
			}
		case "hayvanseviye":
			kullaniciParasi := int64(paraCek(interaction.Member.User.ID))
			var seviye int = hayvanSeviye(interaction.Member.User.ID)
			kalanPara := int(kullaniciParasi) - (seviye * 5000)
			if seviye*5000 <= int(kullaniciParasi) {
				paraKayit(session, interaction.Member.User.ID, int64(kalanPara))
				hayvanGuncelle(interaction.Member.User.ID)
				mesaj := fmt.Sprintf("Hayvanlarınızın, %s MC karşılığında seviyesi yükseltildi, şu an %d seviye hayvanlara sahipsiniz. :star2:", formatNumber(hayvanSeviye(interaction.Member.User.ID)*5000), hayvanSeviye(interaction.Member.User.ID))
				embed := &discordgo.MessageEmbed{
					Description: mesaj,
					Timestamp:   currentTime.Format(time.RFC3339),
				}
				embedGonder(session, interaction, embed)
			} else {
				embed := &discordgo.MessageEmbed{
					Description: fmt.Sprintf("Yeterli paranız yok. Toplam fiyat: %s MC :confused:", formatNumber(hayvanSeviye(interaction.Member.User.ID)*5000)),
				}
				embedGonder(session, interaction, embed)
			}
		case "tohumfiyatları":
			embed := &discordgo.MessageEmbed{
				Title:       "Tohum Fiyatları",
				Description: "Buğday Tohumu: `5.000 MC`\nHavuç Tohumu: `3.500 MC`",
				Timestamp:   currentTime.Format(time.RFC3339),
			}
			embedGonder(session, interaction, embed)
		case "mahsülfiyatları":
			embed := &discordgo.MessageEmbed{
				Title:       "Mahsül Fiyatları",
				Description: "Buğday: `8.500 MC`\nHavuç: `6.000 MC`",
				Timestamp:   currentTime.Format(time.RFC3339),
			}
			embedGonder(session, interaction, embed)
		case "mahsülsat":
			mahsul := interaction.ApplicationCommandData().Options[0].StringValue()
			adet := interaction.ApplicationCommandData().Options[1].IntValue()
			if strings.Contains(mahsul, "Buğday") {
				if int(adet) > mahsulSayi(1, interaction.Member.User.ID) {
					embed := &discordgo.MessageEmbed{
						Description: "Çiftliğinizde bu kadar Buğday bulunmuyor. :confused:",
					}
					embedGonder(session, interaction, embed)
					return
				}
			} else if strings.Contains(mahsul, "Havuç") {
				if int(adet) > mahsulSayi(2, interaction.Member.User.ID) {
					embed := &discordgo.MessageEmbed{
						Description: "Çiftliğinizde bu kadar Havuç bulunmuyor. :confused:",
					}
					embedGonder(session, interaction, embed)
					return
				}
			} else {
				embed := &discordgo.MessageEmbed{
					Description: "Geçersiz mahsül türü. :confused:",
				}
				embedGonder(session, interaction, embed)
				return
			}

			mahsulSat(mahsul, adet, interaction.Member.User.ID)
			var mesaj string
			if strings.Contains(mahsul, "Buğday") {
				mesaj = fmt.Sprintf("%d adet %s sattınız! Aldığınız para: %s MC :slight_smile:", adet, mahsul, formatNumber(int(adet*8500)))
			} else if strings.Contains(mahsul, "Havuç") {
				mesaj = fmt.Sprintf("%d adet %s sattınız! Aldığınız para: %s MC :slight_smile:", adet, mahsul, formatNumber(int(adet*6000)))
			}
			embed := &discordgo.MessageEmbed{
				Description: mesaj,
			}
			embedGonder(session, interaction, embed)
		case "tohumek":
			tohum := interaction.ApplicationCommandData().Options[0].StringValue()
			adet := interaction.ApplicationCommandData().Options[1].IntValue()
			if strings.Contains(tohum, "Buğday Tohumu") {
				if int(adet)+get_EkiliTohum(1, interaction.Member.User.ID) > ciftlikSeviye(interaction.Member.User.ID)*25 || int(adet)+mahsulSayi(1, interaction.Member.User.ID) > ciftlikSeviye(interaction.Member.User.ID)*25 {
					embed := &discordgo.MessageEmbed{
						Description: "Buğday stoğunuz bu kadar mahsülü almıyor. :confused:",
					}
					embedGonder(session, interaction, embed)
					return
				}
				if adet > get_Tohum(tohum, interaction.Member.User.ID) {
					embed := &discordgo.MessageEmbed{
						Description: "Bu kadar Buğday tohumunuz yok. :confused:",
					}
					embedGonder(session, interaction, embed)
					return
				}
			} else if strings.Contains(tohum, "Havuç Tohumu") {
				if int(adet)+get_EkiliTohum(2, interaction.Member.User.ID) > ciftlikSeviye(interaction.Member.User.ID)*25 || int(adet)+mahsulSayi(2, interaction.Member.User.ID) > ciftlikSeviye(interaction.Member.User.ID)*25 {
					embed := &discordgo.MessageEmbed{
						Description: "Havuç stoğunuz bu kadar mahsülü almıyor. :confused:",
					}
					embedGonder(session, interaction, embed)
					return
				}
				if adet > get_Tohum(tohum, interaction.Member.User.ID) {
					embed := &discordgo.MessageEmbed{
						Description: "Bu kadar Havuç tohumunuz yok. :confused:",
					}
					embedGonder(session, interaction, embed)
					return
				}
			} else {
				embed := &discordgo.MessageEmbed{
					Description: "Geçersiz tohum türü. :confused:",
				}
				embedGonder(session, interaction, embed)
				return
			}

			tohumEk(tohum, adet, interaction.Member.User.ID)
			embed := &discordgo.MessageEmbed{
				Description: fmt.Sprintf("%d adet %s ektiniz! :slight_smile:", adet, tohum),
			}
			embedGonder(session, interaction, embed)
		case "tohumsatinal":
			tohum := interaction.ApplicationCommandData().Options[0].StringValue()
			adet := interaction.ApplicationCommandData().Options[1].IntValue()
			var fiyat int64
			if strings.Contains(tohum, "Buğday Tohumu") {
				fiyat = 5000
			} else if strings.Contains(tohum, "Havuç Tohumu") {
				fiyat = 3500
			} else {
				embed := &discordgo.MessageEmbed{
					Description: "Geçersiz tohum türü.",
				}
				embedGonder(session, interaction, embed)
				return
			}

			toplam := adet * fiyat
			kullaniciParasi := int64(paraCek(interaction.Member.User.ID))
			kalanPara := kullaniciParasi - toplam
			if toplam <= int64(kullaniciParasi) {
				paraKayit(session, interaction.Member.User.ID, kalanPara)
				tohumVer(tohum, adet, interaction.Member.User.ID)
				embed := &discordgo.MessageEmbed{
					Description: fmt.Sprintf("%d adet %s satın alındı. Toplam fiyat: %s MC :slight_smile:", adet, tohum, formatNumber(int(toplam))),
				}
				embedGonder(session, interaction, embed)
			} else {
				embed := &discordgo.MessageEmbed{
					Description: fmt.Sprintf("Yeterli paranız yok. Toplam fiyat: %s MC :confused:", formatNumber(int(toplam))),
				}
				embedGonder(session, interaction, embed)
			}
		case "hayvanal":
			hayvan := interaction.ApplicationCommandData().Options[0].StringValue()
			adet := interaction.ApplicationCommandData().Options[1].IntValue()
			var fiyat int64
			if strings.Contains(hayvan, "İnek") {
				fiyat = 25000
			} else if strings.Contains(hayvan, "Koyun") {
				fiyat = 12500
			} else if strings.Contains(hayvan, "Tavuk") {
				fiyat = 5000
			} else {
				embed := &discordgo.MessageEmbed{
					Description: "Geçersiz hayvan türü.",
				}
				embedGonder(session, interaction, embed)
				return
			}

			toplam := adet * fiyat
			kullaniciParasi := int64(paraCek(interaction.Member.User.ID))
			kalanPara := kullaniciParasi - toplam
			if toplam <= int64(kullaniciParasi) {
				paraKayit(session, interaction.Member.User.ID, kalanPara)
				hayvanOlustur(hayvan, adet, interaction.Member.User.ID)
				embed := &discordgo.MessageEmbed{
					Description: fmt.Sprintf("%d adet %s satın alındı. Toplam fiyat: %s MC :slight_smile:", adet, hayvan, formatNumber(int(toplam))),
				}
				embedGonder(session, interaction, embed)
			} else {
				embed := &discordgo.MessageEmbed{
					Description: fmt.Sprintf("Yeterli paranız yok. Toplam fiyat: %s MC :confused:", formatNumber(int(toplam))),
				}
				embedGonder(session, interaction, embed)
			}

		case "girişayarla":
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
		case "kanaloluştur":
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
		case "mesajsil":
			messageCountstr := interaction.ApplicationCommandData().Options[0].StringValue()
			messageCount, err := strconv.Atoi(messageCountstr)
			if err != nil {
				return
			}

			messages, err := session.ChannelMessages(interaction.ChannelID, messageCount, "", "", "")
			if err == nil {
				for _, msg := range messages {
					err := session.ChannelMessageDelete(interaction.ChannelID, msg.ID)
					if err != nil {
						continue
					}
					currentTime := time.Now()
					embed := &discordgo.MessageEmbed{
						Description: messageCountstr + " adet mesaj silindi.",
						Timestamp:   currentTime.Format(time.RFC3339),
					}
					embedGonder(session, interaction, embed)
				}
			}
		case "rolseç":
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
