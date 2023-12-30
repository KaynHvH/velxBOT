package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var dogs = []string{
	// rewrite that to fetch data from some api
}

var answers = []string{
	"Yes", "Of course", "No", "Maybe", "I don't know", "Probably yes", "Probably not", "I don't think so",
}

func (h *Handler) HandleFunCommands(m *discordgo.MessageCreate) {
	if m.Author.ID == h.s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")
	if args[0] != PREFIX {
		return
	}

	switch args[1] {
	case "ping":
		fmt.Println("ping comamnd recognized. Args:", args)
		h.pingCommand(m)
	case "dog":
		fmt.Println("dog  comamnd recognized. Args:", args)
		h.sendRandomDogImage(m)
	case "answer":
		fmt.Println("asnwer comamnd recognized. Args:", args)
		h.sendRandomAnswer(m, args)
	case "whois":
		fmt.Println("whois comamnd recognized. Args:", args)
		h.whoisCommand(m, args)
	}
}

func (h *Handler) pingCommand(m *discordgo.MessageCreate) {
	h.s.ChannelMessageSend(m.ChannelID, "sronk")
}

func (h *Handler) sendRandomDogImage(m *discordgo.MessageCreate) {
	randIndex := rand.Intn(len(dogs))
	dog := dogs[randIndex]

	embed := createEmbed("Random dog image", "Requested by: "+m.Author.Username, dog, m.Author.AvatarURL(""), "dog image")
	h.s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func (h *Handler) sendRandomAnswer(m *discordgo.MessageCreate, args []string) {
	randIndex := rand.Intn(len(answers))
	question := strings.Join(args[2:], " ")

	description := "Answer: **" + answers[randIndex] + "**"
	embed := createEmbed("User "+m.Author.Username+" asked: "+question, "Requested by: "+m.Author.Username, "", m.Author.AvatarURL(""), description)
	h.s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func (h *Handler) whoisCommand(m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		h.s.ChannelMessageSend(m.ChannelID, "Please mention a user to check.")
		return
	}
	mention := args[2]
	userID := strings.Trim(mention, "<@!>")
	user, err := h.s.User(userID)
	if err != nil {
		h.s.ChannelMessageSend(m.ChannelID, "Error retrieving user information.")
		return
	}
	fields := []*discordgo.MessageEmbedField{
		{"Username", user.Username, false},
		{"ID", user.ID, false},
		{"MFA", strconv.FormatBool(user.MFAEnabled), false},
		{"Is bot", strconv.FormatBool(user.Bot), false},
	}
	embed := createEmbed("User Information", "Requested by: "+m.Author.Username, user.AvatarURL(""), m.Author.AvatarURL(""), "", fields...)
	h.s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func createEmbed(title, authorName, thumbnailURL, authorIconURL, description string, fields ...*discordgo.MessageEmbedField) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    authorName,
			IconURL: authorIconURL,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: thumbnailURL,
		},
		Fields: fields,
	}
}
