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
		Name: "**Moderasyon**",
	},
	{
		Value: "`/rolsec:` Girişte üyelere verilecek rolü belirlemek için kullanılır.",
	},
	{
		Value: "`/girisayarla:` Yeni bir üye katıldığında gönderilecek mesajı ve kanalı ayarlamak için kullanılır.",
	},
	{
		Value: "`/kanalolustur:` Yeni bir kanal oluşturmak için kullanılır.",
	},
	{
		Value: "`/kanalsil:` Seçilen kanalı silmek için kullanılır.",
	},
}
