package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Token                    string
	defaultMemberPermissions int64 = discordgo.PermissionManageServer
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

		err = s.GuildMemberRoleAdd(event.GuildID, event.User.ID, rolID)
		if err != nil {
			return
		}
	})
	defer db.Close()

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsGuilds

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
			Name:                     "rolsec",
			Description:              "Bir rol seçin.",
			Type:                     discordgo.ChatApplicationCommand,
			DefaultMemberPermissions: &defaultMemberPermissions,
		},
	}

	fmt.Println("Marin aktif edildi.")
	err = dg.Open()
	if err != nil {
		fmt.Println("bağlantı hatası,", err)
		return
	}

	dg.UpdateStatusComplex(discordgo.UpdateStatusData{AFK: false, Status: (string(discordgo.StatusIdle))})

	_, err = dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, "", command)
	if err != nil {
		log.Println("Error creating slash command:", err)
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
}
