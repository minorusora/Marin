package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Token                    string
	defaultMemberPermissions int64 = discordgo.PermissionManageServer
	kanalYetkisi             int64 = discordgo.PermissionManageChannels
)

const (
	username = ""
	password = ""
	hostname = ""
	dbname   = ""
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Discord oturumu açılamadı:,", err)
		return
	}

	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("%s, veritaban bağlantı hatası.\n", err)
		return
	}

	fmt.Printf("%s, veritabanına başarıyla bağlanıldı.\n", dbname)

	dg.AddHandler(messageCreate)
	dg.AddHandler(interactionCreate)

	dg.AddHandler(func(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
		var rolID string
		err := db.QueryRow("SELECT rol FROM sunucuveri WHERE sunucu_id=?", event.GuildID).Scan(&rolID)
		if err != nil {
			return
		}

		var kanal string
		var mesaj string
		hata := db.QueryRow("SELECT kanal_id, giris_mesaj FROM sunucuveri WHERE sunucu_id=?", event.GuildID).Scan(&kanal, &mesaj)
		if hata != nil {
			log.Printf("%s", hata)
			return
		}

		currentTime := time.Now()
		embed := &discordgo.MessageEmbed{
			Description: event.User.Mention() + " " + mesaj,
			Timestamp:   currentTime.Format(time.RFC3339),
		}

		_, err = s.ChannelMessageSendEmbed(kanal, embed)
		if err != nil {
			return
		}
		err = s.GuildMemberRoleAdd(event.GuildID, event.User.ID, rolID)
		if err != nil {
			return
		}
	})
	defer db.Close()

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsGuilds | discordgo.IntentGuildPresences | discordgo.IntentsGuildPresences

	command := []*discordgo.ApplicationCommand{
		{
			Name:        "avatar",
			Description: "Marin avatarı gönderir.",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "kişi",
					Description: "Bir kişi seçin.",
					Required:    false,
				},
			},
		},
		{
			Name:        "param",
			Description: "Paranızı gösterir.",
			Type:        discordgo.ChatApplicationCommand,
		},
		{
			Name:        "seviyem",
			Description: "Seviyenizi gösterir.",
			Type:        discordgo.ChatApplicationCommand,
		},
		{
			Name:        "çiftliğim",
			Description: "Çiftliğinizin seviyesini ve sahip olduğunuz hayvanları gösterir.",
			Type:        discordgo.ChatApplicationCommand,
		},
		{
			Name:        "hayvanal",
			Description: "Çiftliğinize hayvan almak için kullanılır.",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "hayvan",
					Description: "Alacağınız hayvanı seçin.",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "İnek",
							Value: "İnek",
						},
						{
							Name:  "Koyun",
							Value: "Koyun",
						},
						{
							Name:  "Tavuk",
							Value: "Tavuk",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "adet",
					Description: "Kaç adet alacağınızı girin.",
					Required:    true,
				},
			},
		},
		{
			Name:                     "rolseç",
			Description:              "Girişte üyelere verilecek rolü seçmek için kullanılır.",
			Type:                     discordgo.ChatApplicationCommand,
			DefaultMemberPermissions: &defaultMemberPermissions,
		},
		{
			Name:        "girişayarla",
			Description: "Bir giriş mesajı ayarlayın.",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "kanal",
					Description: "Bir kanal seçin.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "mesaj",
					Description: "Gönderilecek mesajı girin.",
					Required:    true,
				},
			},
			DefaultMemberPermissions: &defaultMemberPermissions,
		},
		{
			Name:        "kanaloluştur",
			Description: "Kanal oluşturmak için kullanılır.",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "kategori",
					Description: "Bir kategori seçin.",
					Required:    true,
				},

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "isim",
					Description: "Oluşturmak istediğiniz kanalın adını girin.",
					Required:    true,
				},
			},
			DefaultMemberPermissions: &kanalYetkisi,
		},
		{
			Name:        "kanalsil",
			Description: "Kanal silmek için kullanılır.",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "kanal",
					Description: "Bir kanal seçin.",
					Required:    true,
				},
			},
			DefaultMemberPermissions: &kanalYetkisi,
		},
		{
			Name:        "yardım",
			Description: "Komutları gösterir.",
			Type:        discordgo.ChatApplicationCommand,
		},
	}

	fmt.Println("Marin aktif edildi.")
	err = dg.Open()
	if err != nil {
		fmt.Println("bağlantı hatası,", err)
		return
	}
	defer dg.Close()

	go ciftlikTimer()

	activity := discordgo.Activity{
		Name: "/yardım",
		Type: discordgo.ActivityTypeGame,
	}

	dg.UpdateStatusComplex(discordgo.UpdateStatusData{

		Activities: []*discordgo.Activity{
			&activity,
		},
		AFK:    false,
		Status: (string(discordgo.StatusIdle)),
	})

	_, err = dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, "", command)
	if err != nil {
		log.Println("Error creating slash command:", err)
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	xpKontrol(s, m.Author.ID, m.Member.GuildID, m.ChannelID)
	mesaj := strings.ToLower(m.Content)
	if mesaj == "selam" {
		s.ChannelMessageSendReply(m.ChannelID, "Selam!", m.Reference())
	} else if mesaj == "merhaba" {
		s.ChannelMessageSendReply(m.ChannelID, "Merhaba!", m.Reference())
	}
}
