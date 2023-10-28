package main

import (
	"github.com/bwmarrin/discordgo"
)

var yardimEmbed = []*discordgo.MessageEmbedField{
	{
		Name: "**Eğlence**",
	},
	{
		Value: "`/avatar:` Kendinizin veya başkasının avatarını gönderir.",
	},
	{
		Value: "`/param:` Paranıza bakmak için kullanılır.",
	},
	{
		Value: "`/çiftliğim:` Çiftliğinizin durumuna bakmak için kullanılır.",
	},
	{
		Value: "`/hayvanal:` Çiftliğinize hayvan almak için kullanılır.",
	},
	{
		Name: "**Moderasyon**",
	},
	{
		Value: "`/rolseç:` Girişte üyelere verilecek rolü belirlemek için kullanılır.",
	},
	{
		Value: "`/girişayarla:` Yeni bir üye katıldığında gönderilecek mesajı ve kanalı ayarlamak için kullanılır.",
	},
	{
		Value: "`/kanaloluştur:` Yeni bir kanal oluşturmak için kullanılır.",
	},
	{
		Value: "`/kanalsil:` Seçilen kanalı silmek için kullanılır.",
	},
}
