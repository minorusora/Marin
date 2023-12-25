package main

import (
	"github.com/bwmarrin/discordgo"
)

var yardimEmbed = []*discordgo.MessageEmbedField{
	{
		Name: " ",
	},
	{
		Name: "**EĞLENCE**",
	},
	{
		Value: "`/avatar:` Kendinizin veya başkasının avatarını gönderir.",
	},
	{
		Value: "`/seviyem:` Seviyenizi ve XP durumunuzu gösterir.",
	},
	{
		Value: " ",
	},
	{
		Name: "**OYUN**",
	},
	{
		Value: "`/param:` Paranıza bakmak için kullanılır.",
	},
	{
		Value: "`/çiftliğim:` Çiftliğinizin istatistiklerine bakmak için kullanılır.",
	},
	{
		Value: "`/çiftlikseviye:` Çiftliğinizin seviyesini yükseltmek için kullanılır.",
	},
	{
		Value: "`/hayvanal:` Çiftliğinize hayvan almak için kullanılır.",
	},
	{
		Value: "`/hayvanfiyatları:` Hayvan fiyatlarına bakmak için kullanılır.",
	},
	{
		Value: "`/hayvanseviye:` Hayvanlarınızın seviyesini yükseltmek için kullanılır.",
	},
	{
		Value: "`/tohumsatinal:` Tohum satın almak için kullanılır.",
	},
	{
		Value: "`/tohumek:` Tohum ekmek için kullanılır.",
	},
	{
		Value: "`/tohumfiyatları:` Tohum fiyatlarına bakmak için kullanılır.",
	},
	{
		Value: "`/mahsülfiyatları:` Mahsül satış fiyatlarına bakmak için kullanılır.",
	},
	{
		Value: "`/mahsülsat:` Mahsül satmak için kullanılır.",
	},
	{
		Value: " ",
	},
	{
		Name: "**MODERASYON**",
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
	{
		Value: "`/mesajsil:` Kanaldaki mesajları silmek için kullanılır.",
	},
}
