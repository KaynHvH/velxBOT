package moderation

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strings"
	"sync"
)

var muteRolesMutex sync.Mutex
var muteRoles = make(map[string]string)

func hasPermission(s *discordgo.Session, guildID, authorID string, permission int64) bool {
	member, err := s.GuildMember(guildID, authorID)
	if err != nil {
		log.Println("Error retrieving guild member information:", err)
		return false
	}

	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			log.Println("Error retrieving role information:", err)
			continue
		}

		if (role.Permissions & permission) == permission {
			return true
		}
	}

	return false
}

func HandleModerationCommands(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")

	if len(args) == 0 || args[0] != os.Getenv("PREFIX") {
		return
	}

	if len(args) >= 2 {
		switch args[1] {
		case "ban":
			fmt.Println("Ban command recognized. Args:", args)
			banUser(s, m, args)
		case "kick":
			fmt.Println("Kick command recognized. Args:", args)
			kickUser(s, m, args)
		case "mute":
			fmt.Println("Mute command recognized. Args:", args)
			muteUser(s, m, args)
		case "unmute":
			fmt.Println("Unmute command recognized. Args:", args)
			unmuteUser(s, m, args)
		case "poll":
			fmt.Println("Poll command recognized. Args:", args)
			poll(s, m, args)
		case "nick", "nickname":
			fmt.Println("Nickname command recognized. Args:", args)
			changeUserNickname(s, m, args)
		}
	}
}

func poll(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please enter poll content.")
		return
	}

	if hasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator) {
		content := strings.Join(args[2:len(args)], " ")

		embed := discordgo.MessageEmbed{
			Title:       ",," + content + "''",
			Description: "React to the poll down below ",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Requested by: " + m.Author.Username,
				IconURL: m.Author.AvatarURL(""),
			},
		}
		msg, err := s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}

		// Reactions to the embed ;)
		err = s.MessageReactionAdd(m.ChannelID, msg.ID, "👍")
		if err != nil {
			log.Println("Error adding thumbs up reaction:", err)
		}

		err = s.MessageReactionAdd(m.ChannelID, msg.ID, "👎")
		if err != nil {
			log.Println("Error adding thumbs down reaction:", err)
		}

		return
	} else {
		s.ChannelMessageSend(m.ChannelID, "You don't have permissions to create a query!")
	}
}

func banUser(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Please mention a user to ban.")
		return
	}

	userID := strings.Trim(args[2], "<@!>")
	user, err := s.User(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error retrieving user information.")
		return
	}

	if hasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator) {
		reason := strings.Join(args[3:len(args)], " ")

		err = s.GuildBanCreate(m.GuildID, user.ID, 0)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Can't ban the user")
			log.Printf("Error banning user: %s", err)
			return
		}

		embed := discordgo.MessageEmbed{
			Title:       "User " + user.Username + " has been banned!",
			Description: "Reason: " + reason,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Requested by: " + m.Author.Username,
				IconURL: m.Author.AvatarURL(""),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		return
	} else {
		embed := discordgo.MessageEmbed{
			Title:       "User " + m.Author.Username + " tried to use administrator command",
			Description: "You don't have specific permissions to ban someone ;)",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Requested by: " + m.Author.Username,
				IconURL: m.Author.AvatarURL(""),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	}
	return
}

func kickUser(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Please mention a user to kick.")
		return
	}

	userID := strings.Trim(args[2], "<@!>")
	user, err := s.User(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error retrieving user information.")
		return
	}

	if hasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator) {
		reason := strings.Join(args[3:len(args)], " ")

		err = s.GuildMemberDeleteWithReason(m.GuildID, user.ID, reason)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Can't kick the user")
			log.Printf("Error kicking user: %s", err)
			return
		}

		embed := discordgo.MessageEmbed{
			Title:       "User " + user.Username + " has been kicked!",
			Description: "Reason: " + reason,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Requested by: " + m.Author.Username,
				IconURL: m.Author.AvatarURL(""),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	} else {
		embed := discordgo.MessageEmbed{
			Title:       "User " + m.Author.Username + " tried to use administrator command",
			Description: "You don't have specific permissions to kick someone ;)",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Requested by: " + m.Author.Username,
				IconURL: m.Author.AvatarURL(""),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	}
	return
}

func muteUser(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Please mention a user to mute.")
		return
	}

	userID := strings.Trim(args[2], "<@!>")
	user, err := s.User(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error retrieving user information.")
		return
	}

	if hasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator) {
		reason := strings.Join(args[3:len(args)], " ")
		muteRoleID, err := getMuteRoleID(s, m.GuildID)

		err = s.GuildMemberRoleAdd(m.GuildID, user.ID, muteRoleID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Can't mute the user")
			log.Printf("Error muting user: %s", err)
			return
		}

		embed := discordgo.MessageEmbed{
			Title:       "User " + user.Username + " has been muted!",
			Description: "Reason: " + reason,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Requested by: " + m.Author.Username,
				IconURL: m.Author.AvatarURL(""),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		return
	} else {
		embed := discordgo.MessageEmbed{
			Title:       "User " + m.Author.Username + " tried to use administrator command",
			Description: "You don't have specific permissions to mute someone ;)",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Requested by: " + m.Author.Username,
				IconURL: m.Author.AvatarURL(""),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		return
	}
}

func unmuteUser(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Please mention a user to unmute.")
		return
	}

	userID := strings.Trim(args[2], "<@!>")
	user, err := s.User(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error retrieving user information.")
		return
	}

	if hasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator) {
		muteRoleID, err := getMuteRoleID(s, m.GuildID)

		err = s.GuildMemberRoleRemove(m.GuildID, user.ID, muteRoleID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Can't unmute the user")
			log.Printf("Error muting user: %s", err)
			return
		}

		embed := discordgo.MessageEmbed{
			Title: "User " + user.Username + " has been unmuted!",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Requested by: " + m.Author.Username,
				IconURL: m.Author.AvatarURL(""),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		return
	} else {
		embed := discordgo.MessageEmbed{
			Title:       "User " + m.Author.Username + " tried to use administrator command",
			Description: "You don't have specific permissions to unmute someone ;)",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Requested by: " + m.Author.Username,
				IconURL: m.Author.AvatarURL(""),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		return
	}
}

func getMuteRoleID(s *discordgo.Session, guildID string) (string, error) {
	muteRolesMutex.Lock()
	defer muteRolesMutex.Unlock()

	if roleID, ok := muteRoles[guildID]; ok {
		return roleID, nil
	}

	muteRoleID := "921805059581423657"

	muteRoles[guildID] = muteRoleID

	return muteRoleID, nil
}

func changeUserNickname(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Please mention a user to change the nickname.")
		return
	}

	userID := strings.Trim(args[2], "<@!>")
	user, err := s.User(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error retrieving user information.")
		return
	}

	if hasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator) {
		newNickname := strings.Join(args[3:len(args)], " ")

		err = s.GuildMemberNickname(m.GuildID, user.ID, newNickname)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Can't change user nickname")
			log.Printf("Error changing user nickname: %s", err)
			return
		}

		embed := discordgo.MessageEmbed{
			Title: user.Username + " nickname has been changed to **" + newNickname + "**",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Requested by: " + m.Author.Username,
				IconURL: m.Author.AvatarURL(""),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		return
	} else {
		embed := discordgo.MessageEmbed{
			Title:       "User " + m.Author.Username + " tried to use administrator command",
			Description: "You don't have specific permissions to change someone's nickname ;)",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Requested by: " + m.Author.Username,
				IconURL: m.Author.AvatarURL(""),
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		return
	}
}
