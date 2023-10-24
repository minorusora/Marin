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
	Token string
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
	defer db.Close()

	fmt.Printf("%s, veritabanına başarıyla bağlanıldı.\n", dbname)

	dg.AddHandler(messageCreate)
	dg.AddHandler(interactionHandler)

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

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsGuilds

	fmt.Println("Marin aktif edildi.")
	dg.UpdateStatusComplex(discordgo.UpdateStatusData{AFK: false, Status: string(discordgo.StatusIdle)})

	err = dg.Open()
	if err != nil {
		fmt.Println("bağlantı hatası,", err)
		return
	}

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "avatar",
			Description: "Marin avatarınızı gönderir.",
		},
		{
			Name:        "rolayarla",
			Description: "Sunucuya yeni bir üye katıldığında verilecek rolü ayarlamak için kullanılır.",
		},
	}

	dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, "1165767884916658307", commands)

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
