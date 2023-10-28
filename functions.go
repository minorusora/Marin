package main

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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

func formatNumber(num int) string {
	numStr := strconv.Itoa(num)

	var formattedNum string
	for i, digit := range numStr {
		if i > 0 && (len(numStr)-i)%3 == 0 {
			formattedNum += "."
		}
		formattedNum += string(digit)
	}

	return formattedNum
}

func paraKayit(session *discordgo.Session, userID string, para int64) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		_, err := db.Exec("UPDATE kullaniciveri SET para = ? WHERE kisi_id = ?", para, userID)
		if err != nil {
			panic(err.Error())
		}
	} else {
		_, err := db.Exec("insert into kullaniciveri(kisi_id, para, ciftlik_seviye) values (?, ?, ?)", userID, 1250, 1)
		if err != nil {
			panic(err.Error())
		}
	}

	defer db.Close()
}

func paraCek(userID string) int {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var para int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&count)
	if err != nil {
		panic(err.Error())
	}

	if count <= 0 {
		_, err := db.Exec("INSERT INTO kullaniciveri (kisi_id, para, ciftlik_seviye) VALUES (?, ?, ?)", userID, 1250, 1)
		if err != nil {
			panic(err.Error())
		} else {
			return para
		}
	} else {
		err = db.QueryRow("SELECT para FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&para)
		if err != nil {
			panic(err.Error())
		}
	}
	return para
}

func inekGet_Count(userID string) int {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var inek_sayi int
	err = db.QueryRow("SELECT COUNT(*) FROM inekler WHERE sahip_id = ?", userID).Scan(&inek_sayi)
	if err != nil {
		panic(err.Error())
	}
	return inek_sayi
}

func koyunGet_Count(userID string) int {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var koyun_sayi int
	err = db.QueryRow("SELECT COUNT(*) FROM koyunlar WHERE sahip_id = ?", userID).Scan(&koyun_sayi)
	if err != nil {
		panic(err.Error())
	}
	return koyun_sayi
}

func tavukGet_Count(userID string) int {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var tavuk_sayi int
	err = db.QueryRow("SELECT COUNT(*) FROM tavuklar WHERE sahip_id = ?", userID).Scan(&tavuk_sayi)
	if err != nil {
		panic(err.Error())
	}
	return tavuk_sayi
}

func ciftlikSeviye(userID string) int {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var ciftlikSeviye int
	err = db.QueryRow("SELECT ciftlik_seviye FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&ciftlikSeviye)
	if err != nil {
		panic(err.Error())
	}
	return ciftlikSeviye
}

func hayvanOlustur(hayvan string, adet int64, userID string) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	for i := 0; i < int(adet); i++ {
		if strings.Contains(hayvan, "İnek") {
			_, err := db.Exec("INSERT INTO inekler (sahip_id, gelir) VALUES (?, ?)", userID, ciftlikSeviye(userID)*10)
			if err != nil {
				return
			}
		} else if strings.Contains(hayvan, "Koyun") {
			_, err = db.Exec("INSERT INTO koyunlar (sahip_id, gelir) VALUES (?, ?)", userID, ciftlikSeviye(userID)*6)
			if err != nil {
				return
			}
		} else if strings.Contains(hayvan, "Tavuk") {
			_, err = db.Exec("INSERT INTO tavuklar (sahip_id, gelir) VALUES (?, ?)", userID, ciftlikSeviye(userID)*3)
			if err != nil {
				return
			}
		}
	}
}
