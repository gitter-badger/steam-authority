package main

import (
	"net/http"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/chi"
	"github.com/steam-authority/steam-authority/websockets"
)

const (
	GUILD_ID   = "407493776597057538"
	GENERAL_ID = "407493777058693121"
)

var (
	discord *discordgo.Session
)

func init() {

	var err error

	// Get client
	discord, err = discordgo.New("Bot NDA1MDQ4MTc1OTI2MzEyOTcx.DVD2tQ.ZsyCsJ6jjYE4Hw6QNP58LWz3GqA")
	if err != nil {
		logger.Error(err)
	}

	// Add websocket listener
	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if !m.Author.Bot {
			websockets.Send(websockets.CHAT, m)
		}
	})

	// Open connection
	err = discord.Open()
	if err != nil {
		logger.Error(err)
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {

	// Get ID from URL
	id := chi.URLParam(r, "id")
	if id == "" {
		id = GENERAL_ID
	}

	// Get channels
	channelsResponse, err := discord.GuildChannels(GUILD_ID)
	if err != nil {
		logger.Error(err)
	}

	channels := make([]*discordgo.Channel, 0)
	for _, v := range channelsResponse {
		if v.Type == discordgo.ChannelTypeGuildText {
			channels = append(channels, v)
		}
	}

	// Get messages
	messagesResponse, err := discord.ChannelMessages(id, 10, "", "", "")
	if err != nil {
		logger.Error(err)
	}

	messages := make([]*discordgo.Message, 0)
	for _, v := range messagesResponse {
		if !v.Author.Bot {
			messages = append(messages, v)
		}
	}

	// Template
	template := chatTemplate{}
	template.Channels = channels
	template.Messages = messages
	template.ChannelID = id

	returnTemplate(w, "chat", template)
}

type chatTemplate struct {
	GlobalTemplate
	Channels  []*discordgo.Channel
	Messages  []*discordgo.Message
	ChannelID string
}
