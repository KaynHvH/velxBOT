package help

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
)

func HandleHelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")

	if len(args) == 0 || args[0] != os.Getenv("PREFIX") {
		return
	}

	if len(args) >= 2 {
		switch args[1] {
		case "help":
			fmt.Println("Help command recognized. Args:", args)

			embed := discordgo.MessageEmbed{
				Title: "Velx Bot commands",
				Description: "**velx dog** - Sends random dog picture by API\n" +
					"**velx answer <question>** - Answers to the question\n" +
					"**velx avatar <@user>** - Displays user's avatar\n" +
					"**velx whois <@user>** - Tells about user\n" +
					"**velx poll <content>** - Creates poll with reactions\n" +
					"**velx dice** - Rolls the dice\n" +
					"**velx ban <@user> <reason>** - Bans user with reason\n" +
					"**velx kick <@user> <reason>** - Kicks user with reason\n" +
					"**velx mute <@user> <reason>** - Mutes user with reason\n" +
					"**velx unmute <@user>** - Unmutes user\n" +
					"**velx nick/nickname <@user> <nickname>** - Changes user nickname\n",
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
