package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

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
		_, err := db.Exec("insert into kullaniciveri(kisi_id, para, ciftlik_seviye) values (?, ?, ?)", userID, 100000, 1)
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
		_, err := db.Exec("INSERT INTO kullaniciveri (kisi_id, para, ciftlik_seviye) VALUES (?, ?, ?)", userID, 100000, 1)
		if err != nil {
			panic(err.Error())
		} else {
			return 100000
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
	err = db.QueryRow("SELECT inek_adet FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&inek_sayi)
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
	err = db.QueryRow("SELECT koyun_adet FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&koyun_sayi)
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
	err = db.QueryRow("SELECT tavuk_adet FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&tavuk_sayi)
	if err != nil {
		panic(err.Error())
	}
	return tavuk_sayi
}

func hayvanSeviye(userID string) int {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var hayvanSeviye int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		err = db.QueryRow("SELECT hayvan_seviye FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&hayvanSeviye)
		if err != nil {
			panic(err.Error())
		}
	} else {
		_, err := db.Exec("INSERT INTO kullaniciveri (kisi_id, para) VALUES (?, ?)", userID, 100000)
		if err != nil {
			panic(err.Error())
		}
	}
	return hayvanSeviye
}

func hayvanGuncelle(userID string) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		_, err := db.Exec("UPDATE kullaniciveri SET hayvan_seviye = hayvan_seviye + 1 WHERE kisi_id = ?", userID)
		if err != nil {
			panic(err.Error())
		}
		_, err = db.Exec("UPDATE kullaniciveri SET inek_gelir = ?, koyun_gelir = ?, tavuk_gelir = ? WHERE kisi_id = ?", hayvanSeviye(userID)*10, hayvanSeviye(userID)*6, hayvanSeviye(userID)*3, userID)
		if err != nil {
			panic(err.Error())
		}
	} else {
		_, err := db.Exec("INSERT INTO kullaniciveri (kisi_id, para) VALUES (?, ?)", userID, 100000)
		if err != nil {
			panic(err.Error())
		}
	}
}

func tohumSayisi(tur int, userID string) int {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var count int
	var sayi int
	err = db.QueryRow("SELECT COUNT(*) FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		if tur == 1 {
			err = db.QueryRow("SELECT bugdaytohum_adet FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&sayi)
			if err != nil {
				panic(err.Error())
			}
		} else if tur == 2 {
			err = db.QueryRow("SELECT havuctohum_adet FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&sayi)
			if err != nil {
				panic(err.Error())
			}
		}
	} else {
		_, err := db.Exec("INSERT INTO kullaniciveri (kisi_id, para) VALUES (?, ?)", userID, 100000)
		if err != nil {
			panic(err.Error())
		}
	}
	return sayi
}

func get_EkiliTohum(tur int, userID string) int {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var count int
	var sayi int
	err = db.QueryRow("SELECT COUNT(*) FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		if tur == 1 {
			err = db.QueryRow("SELECT bugdayekili_adet FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&sayi)
			if err != nil {
				panic(err.Error())
			}
		} else if tur == 2 {
			err = db.QueryRow("SELECT havucekili_adet FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&sayi)
			if err != nil {
				panic(err.Error())
			}
		}
	} else {
		_, err := db.Exec("INSERT INTO kullaniciveri (kisi_id, para) VALUES (?, ?)", userID, 100000)
		if err != nil {
			panic(err.Error())
		}
	}
	return sayi
}

func mahsulSayi(tur int, userID string) int {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var count int
	var sayi int
	err = db.QueryRow("SELECT COUNT(*) FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		if tur == 1 {
			err = db.QueryRow("SELECT bugday_adet FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&sayi)
			if err != nil {
				panic(err.Error())
			}
		} else if tur == 2 {
			err = db.QueryRow("SELECT havuc_adet FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&sayi)
			if err != nil {
				panic(err.Error())
			}
		}
	} else {
		_, err := db.Exec("INSERT INTO kullaniciveri (kisi_id, para) VALUES (?, ?)", userID, 100000)
		if err != nil {
			panic(err.Error())
		}
	}
	return sayi
}

func ciftlikSeviye(userID string) int {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var ciftlikSeviye int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		err = db.QueryRow("SELECT ciftlik_seviye FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&ciftlikSeviye)
		if err != nil {
			panic(err.Error())
		}
	} else {
		_, err := db.Exec("INSERT INTO kullaniciveri (kisi_id, para) VALUES (?, ?)", userID, 100000)
		if err != nil {
			panic(err.Error())
		}
	}
	return ciftlikSeviye
}

func get_Tohum(tohum string, userID string) int64 {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var adet int64
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		if strings.Contains(tohum, "Buğday Tohumu") {
			err = db.QueryRow("SELECT bugdaytohum_adet FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&adet)
			if err != nil {
				panic(err.Error())
			}
		} else if strings.Contains(tohum, "Havuç Tohumu") {
			err = db.QueryRow("SELECT havuctohum_adet FROM kullaniciveri WHERE kisi_id = ?", userID).Scan(&adet)
			if err != nil {
				panic(err.Error())
			}
		}
	} else {
		_, err := db.Exec("INSERT INTO kullaniciveri (kisi_id, para) VALUES (?, ?)", userID, 100000)
		if err != nil {
			panic(err.Error())
		}
	}
	return adet
}

func mahsulSat(mahsul string, adet int64, userID string) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if strings.Contains(mahsul, "Buğday") {
		_, err := db.Exec("UPDATE kullaniciveri SET bugday_adet = bugday_adet - ? WHERE kisi_id = ?", adet, userID)
		if err != nil {
			panic(err.Error())
		}

		kullaniciGuncelle(db, userID, int(8500*adet))

	} else if strings.Contains(mahsul, "Havuç") {
		_, err := db.Exec("UPDATE kullaniciveri SET havuc_adet = havuc_adet - ? WHERE kisi_id = ?", adet, userID)
		if err != nil {
			panic(err.Error())
		}

		kullaniciGuncelle(db, userID, int(6000*adet))
	}
}

func tohumEk(tohum string, adet int64, userID string) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if strings.Contains(tohum, "Buğday Tohumu") {
		_, err := db.Exec("UPDATE kullaniciveri SET bugdayekili_adet = bugdayekili_adet + ? WHERE kisi_id = ?", adet, userID)
		if err != nil {
			panic(err.Error())
		}
		_, err = db.Exec("UPDATE kullaniciveri SET bugdaytohum_adet = bugdaytohum_adet - ? WHERE kisi_id = ?", adet, userID)
		if err != nil {
			panic(err.Error())
		}

	} else if strings.Contains(tohum, "Havuç Tohumu") {
		_, err := db.Exec("UPDATE kullaniciveri SET havucekili_adet = havucekili_adet + ? WHERE kisi_id = ?", adet, userID)
		if err != nil {
			panic(err.Error())
		}

		_, err = db.Exec("UPDATE kullaniciveri SET havuctohum_adet = havuctohum_adet - ? WHERE kisi_id = ?", adet, userID)
		if err != nil {
			panic(err.Error())
		}
	}
}

func tohumVer(tohum string, adet int64, userID string) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if strings.Contains(tohum, "Buğday Tohumu") {
		_, err := db.Exec("UPDATE kullaniciveri SET bugdaytohum_adet = bugdaytohum_adet + ? WHERE kisi_id = ?", adet, userID)
		if err != nil {
			panic(err.Error())
		}
	} else if strings.Contains(tohum, "Havuç Tohumu") {
		_, err := db.Exec("UPDATE kullaniciveri SET havuctohum_adet = havuctohum_adet + ? WHERE kisi_id = ?", adet, userID)
		if err != nil {
			panic(err.Error())
		}
	}
}

func hayvanOlustur(hayvan string, adet int64, userID string) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if strings.Contains(hayvan, "İnek") {
		_, err := db.Exec("UPDATE kullaniciveri SET inek_adet = inek_adet + ?, inek_gelir = ? WHERE kisi_id = ?", adet, hayvanSeviye(userID)*10, userID)
		if err != nil {
			panic(err.Error())
		}
	} else if strings.Contains(hayvan, "Koyun") {
		_, err = db.Exec("UPDATE kullaniciveri SET koyun_adet = koyun_adet + ?, koyun_gelir = ? WHERE kisi_id = ?", adet, hayvanSeviye(userID)*6, userID)
		if err != nil {
			panic(err.Error())
		}
	} else if strings.Contains(hayvan, "Tavuk") {
		_, err = db.Exec("UPDATE kullaniciveri SET tavuk_adet = tavuk_adet + ?, tavuk_gelir = ? WHERE kisi_id = ?", adet, hayvanSeviye(userID)*3, userID)
		if err != nil {
			panic(err.Error())
		}
	}
}

func ciftlikSeviyeYukselt(userID string) {
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
		_, err := db.Exec("UPDATE kullaniciveri SET ciftlik_seviye = ciftlik_seviye + 1 WHERE kisi_id = ?", userID)
		if err != nil {
			panic(err.Error())
		}
	} else {
		_, err := db.Exec("INSERT INTO kullaniciveri (kisi_id, para) VALUES (?, ?)", userID, 100000)
		if err != nil {
			panic(err.Error())
		}
	}
}

func xpCheck(userID string) int {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var xp int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM xp WHERE kisi_id = ?", userID).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		err = db.QueryRow("SELECT xp FROM xp WHERE kisi_id = ?", userID).Scan(&xp)
		if err != nil {
			panic(err.Error())
		}
	} else {
		return 0
	}
	return xp
}

func levelKontrol(userID string) int {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var level int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM xp WHERE kisi_id = ?", userID).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		err = db.QueryRow("SELECT level FROM xp WHERE kisi_id = ?", userID).Scan(&level)
		if err != nil {
			panic(err.Error())
		}
	} else {
		return 1
	}
	return level
}

func xpKontrol(session *discordgo.Session, userID string, guildID string, channelID string) error {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM xp WHERE kisi_id = ? and sunucu_id = ?", userID, guildID).Scan(&count)
	if err != nil {
		panic(err.Error())
	}

	if count <= 0 {
		_, err := db.Exec("INSERT INTO xp (kisi_id, sunucu_id) VALUES (?, ?)", userID, guildID)
		if err != nil {
			panic(err.Error())
		}
	} else {
		toplam := levelKontrol(userID) * 240
		if xpCheck(userID) >= toplam {
			_, err := db.Exec("UPDATE xp SET level = ? WHERE kisi_id = ? and sunucu_id = ?", levelKontrol(userID)+1, userID, guildID)
			if err != nil {
				panic(err.Error())
			}

			mesaj := fmt.Sprintf("<@"+userID+">"+" tebrikler, %d. seviyeye ulaştınız!", levelKontrol(userID))
			session.ChannelMessageSend(channelID, mesaj)
		} else {
			_, err := db.Exec("UPDATE xp SET xp = ? WHERE kisi_id = ? and sunucu_id = ?", xpCheck(userID)+10, userID, guildID)
			if err != nil {
				panic(err.Error())
			}
		}
	}
	return err
}

func inekGelir(db *sql.DB) {
	rows, err := db.Query("SELECT inek_adet, inek_gelir, kisi_id FROM kullaniciveri WHERE kisi_id != 0 and inek_adet != 0")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var gelir int
		var inek_adet int
		var userID string
		err := rows.Scan(&inek_adet, &gelir, &userID)
		if err != nil {
			return
		}
		kullaniciGuncelle(db, userID, gelir*inek_adet)
	}
}

func koyunGelir(db *sql.DB) {
	rows, err := db.Query("SELECT koyun_adet, koyun_gelir, kisi_id FROM kullaniciveri WHERE kisi_id != 0 and koyun_adet != 0")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var gelir int
		var koyun_adet int
		var userID string
		err := rows.Scan(&koyun_adet, &gelir, &userID)
		if err != nil {
			return
		}
		kullaniciGuncelle(db, userID, gelir*koyun_adet)
	}
}

func tavukGelir(db *sql.DB) {
	rows, err := db.Query("SELECT tavuk_adet, tavuk_gelir, kisi_id FROM kullaniciveri WHERE kisi_id != 0 and tavuk_adet != 0")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var gelir int
		var tavuk_adet int
		var userID string
		err := rows.Scan(&tavuk_adet, &gelir, &userID)
		if err != nil {
			return
		}
		kullaniciGuncelle(db, userID, gelir*tavuk_adet)
	}
}

func bugdayYetisme(db *sql.DB, session *discordgo.Session) {
	rows, err := db.Query("SELECT bugdayekili_adet, kisi_id FROM kullaniciveri WHERE kisi_id != 0 and bugdayekili_adet != 0")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var bugday_adet int
		var userID string
		err := rows.Scan(&bugday_adet, &userID)
		if err != nil {
			return
		} else {
			_, err := db.Exec("UPDATE kullaniciveri SET bugday_adet = bugday_adet + bugdayekili_adet WHERE kisi_id = ?", userID)
			if err != nil {
				panic(err.Error())
			}
			kullaniciGuncelle(db, userID, 10*bugday_adet)
			channel, err := session.UserChannelCreate(userID)
			if err == nil {

				embed := &discordgo.MessageEmbed{
					Description: "Buğday tohumlarınız yetişti! :farmer: :woman_farmer:",
				}

				_, err := db.Exec("UPDATE kullaniciveri SET bugdayekili_adet = 0 WHERE kisi_id = ?", userID)
				if err != nil {
					panic(err.Error())
				}

				_, err = session.ChannelMessageSendEmbed(channel.ID, embed)
				if err != nil {
					panic(err.Error())
				}
			}
		}
	}
}

func havucYetisme(db *sql.DB, session *discordgo.Session) {
	rows, err := db.Query("SELECT havucekili_adet, kisi_id FROM kullaniciveri WHERE kisi_id != 0 and havucekili_adet != 0")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var havucekili_adet int
		var userID string
		err := rows.Scan(&havucekili_adet, &userID)
		if err != nil {
			return
		} else {
			_, err := db.Exec("UPDATE kullaniciveri SET havuc_adet = havuc_adet + havucekili_adet WHERE kisi_id = ?", userID)
			if err != nil {
				panic(err.Error())
			}

			kullaniciGuncelle(db, userID, 10*havucekili_adet)
			channel, err := session.UserChannelCreate(userID)
			if err == nil {
				embed := &discordgo.MessageEmbed{
					Description: "Havuç tohumlarınız yetişti! :farmer: :woman_farmer:",
				}

				_, err := db.Exec("UPDATE kullaniciveri SET havucekili_adet = 0 WHERE kisi_id = ?", userID)
				if err != nil {
					panic(err.Error())
				}

				_, err = session.ChannelMessageSendEmbed(channel.ID, embed)
				if err != nil {
					panic(err.Error())
				}
			}
		}
	}
}

func bugdayTimer(session *discordgo.Session) {
	for {
		db, err := sql.Open("mysql", dsn(dbname))
		if err != nil {
			panic(err.Error())
		}
		bugdayYetisme(db, session)
		time.Sleep(2 * time.Hour)
	}
}

func havucTimer(session *discordgo.Session) {
	for {
		db, err := sql.Open("mysql", dsn(dbname))
		if err != nil {
			panic(err.Error())
		}
		havucYetisme(db, session)
		time.Sleep(1 * time.Hour)
	}
}

func ciftlikTimer() {
	for {
		db, err := sql.Open("mysql", dsn(dbname))
		if err != nil {
			panic(err.Error())
		}
		inekGelir(db)
		koyunGelir(db)
		tavukGelir(db)
		time.Sleep(3 * time.Hour)
	}
}

func kullaniciGuncelle(db *sql.DB, userID string, para int) {
	_, err := db.Exec("UPDATE kullaniciveri SET para = para + ? WHERE kisi_id = ?", para, userID)
	if err != nil {
		return
	}
}
