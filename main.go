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

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Marin aktif edildi.")
	dg.UpdateStatusComplex(discordgo.UpdateStatusData{AFK: false, Status: string(discordgo.StatusIdle)})

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "avatar" {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name: m.Author.Username,
			},
			Image: &discordgo.MessageEmbedImage{
				URL:    m.Author.AvatarURL(m.Author.Avatar),
				Width:  128,
				Height: 128,
			},
		})
	}
}
