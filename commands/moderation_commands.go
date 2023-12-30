package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const PREFIX = "your_prefix_here"

var muteRoles = make(map[string]string)

func (h *Handler) HandleModerationCommands(m *discordgo.MessageCreate) {
	if m.Author.ID == h.s.State.User.ID {
		return
	}

	args := strings.Fields(m.Content)

	if len(args) == 0 || args[0] != PREFIX {
		return
	}

	if len(args) >= 2 {
		switch args[1] {
		case "ban":
			fmt.Println("Ban command recognized. Args:", args)
			h.banKickUser(m, args, discordgo.PermissionBanMembers, "banned")
		case "kick":
			fmt.Println("Kick command recognized. Args:", args)
			h.banKickUser(m, args, discordgo.PermissionKickMembers, "kicked")
		case "mute":
			fmt.Println("Mute command recognized. Args:", args)
			h.muteUser(m, args)
		}
	}
}

func (h *Handler) hasPermission(guildID, authorID string, permission int64) bool {
	member, err := h.s.GuildMember(guildID, authorID)
	if err != nil {
		log.Println("Error retrieving guild member information:", err)
		return false
	}

	for _, roleID := range member.Roles {
		role, err := h.s.State.Role(guildID, roleID)
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

func (h *Handler) banKickUser(m *discordgo.MessageCreate, args []string, permission int64, action string) {
	if len(args) < 3 {
		h.s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Please mention a user to %s.", action))
		return
	}

	userID := strings.Trim(args[2], "<@!>")
	user, err := h.s.User(userID)
	if err != nil {
		h.s.ChannelMessageSend(m.ChannelID, "Error retrieving user information.")
		return
	}

	if h.hasPermission(m.GuildID, m.Author.ID, discordgo.PermissionAdministrator) {
		reason := strings.Join(args[3:], " ")

		err = h.s.GuildBanCreate(m.GuildID, user.ID, 0)
		if err != nil {
			h.s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Can't %s the user", action))
			log.Fatalf("Error %sing user: %s", action, err)
			return
		}

		embed := discordgo.MessageEmbed{
			Title:       fmt.Sprintf("User %s has been %s!", user.Username, action),
			Description: fmt.Sprintf("Reason: %s", reason),
			Author: &discordgo.MessageEmbedAuthor{
				Name:    fmt.Sprintf("Requested by: %s", m.Author.Username),
				IconURL: m.Author.AvatarURL(""),
			},
		}
		h.s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	} else {
		embed := discordgo.MessageEmbed{
			Title:       fmt.Sprintf("User %s tried to use administrator command", m.Author.Username),
			Description: fmt.Sprintf("You don't have specific permissions to %s someone ;)", action),
			Author: &discordgo.MessageEmbedAuthor{
				Name:    fmt.Sprintf("Requested by: %s", m.Author.Username),
				IconURL: m.Author.AvatarURL(""),
			},
		}
		h.s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	}
}

func (h *Handler) muteUser(m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		h.s.ChannelMessageSend(m.ChannelID, "Please mention a user to mute.")
		return
	}

	userID := strings.Trim(args[2], "<@!>")
	user, err := h.s.User(userID)
	if err != nil {
		h.s.ChannelMessageSend(m.ChannelID, "Error retrieving user information.")
		return
	}

	if h.hasPermission(m.GuildID, m.Author.ID, discordgo.PermissionAdministrator) {
		reason := strings.Join(args[3:], " ")
		muteRoleID, err := h.getMuteRoleID(m.GuildID)

		err = h.s.GuildMemberRoleAdd(m.GuildID, user.ID, muteRoleID)
		if err != nil {
			h.s.ChannelMessageSend(m.ChannelID, "Can't mute the user")
			log.Fatalf("Error muting user: %s", err)
			return
		}

		embed := discordgo.MessageEmbed{
			Title:       fmt.Sprintf("User %s has been muted!", user.Username),
			Description: fmt.Sprintf("Reason: %s", reason),
			Author: &discordgo.MessageEmbedAuthor{
				Name:    fmt.Sprintf("Requested by: %s", m.Author.Username),
				IconURL: m.Author.AvatarURL(""),
			},
		}
		h.s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	} else {
		embed := discordgo.MessageEmbed{
			Title:       fmt.Sprintf("User %s tried to use administrator command", m.Author.Username),
			Description: "You don't have specific permissions to ban someone ;)",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    fmt.Sprintf("Requested by: %s", m.Author.Username),
				IconURL: m.Author.AvatarURL(""),
			},
		}
		h.s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	}
}

func (h *Handler) getMuteRoleID(guildID string) (string, error) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if roleID, ok := muteRoles[guildID]; ok {
		return roleID, nil
	}

	muteRoleID := "921805059581423657"

	muteRoles[guildID] = muteRoleID

	return muteRoleID, nil
}
