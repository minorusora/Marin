package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.AddHandler(interactionHandler)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Marin aktif edildi.")
	dg.UpdateGameStatus(0, "/yardım")
	dg.UpdateStatusComplex(discordgo.UpdateStatusData{AFK: false, Status: string(discordgo.StatusIdle)})

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "avatar",
			Description: "Marin avatarınızı gönderir.",
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
