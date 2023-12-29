package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const PREFIX string = "velx"

func HandleHelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")

	if len(args) == 0 || args[0] != PREFIX {
		return
	}

	if len(args) >= 2 {
		if args[1] == "help" {
			fmt.Println("Help command recognized. Args:", args)

			embed := discordgo.MessageEmbed{
				Title: "Velx Bot commands",
				Description: "velx dog - Sends random dog picture\n" +
					"velx answer <question> - Answers to the question\n" +
					"velx whois <@user> - Tells about user\n" +
					"velx ban <@user> <reason> - Bans user with reason\n" +
					"velx kick <@user> <reason> - Kicks user with reason\n" +
					"velx mute <@user> <reason> - Mutes user with reason",
				Author: &discordgo.MessageEmbedAuthor{
					Name:    "Requested by: " + m.Author.Username,
					IconURL: m.Author.AvatarURL(""),
				},
			}

			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
			return
		}
	}
}
